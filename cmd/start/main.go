package main

import (
	"context"
	"errors"
	"fmt"
	permify_grpc "github.com/Permify/permify-go/grpc"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
	"github.com/redis/rueidis/rueidisotel"
	"github.com/rsmidt/soccerbuddy/gen/eventregistry"
	"github.com/rsmidt/soccerbuddy/internal/app/commands"
	"github.com/rsmidt/soccerbuddy/internal/app/queries"
	"github.com/rsmidt/soccerbuddy/internal/config"
	"github.com/rsmidt/soccerbuddy/internal/domain"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/grpc"
	"github.com/rsmidt/soccerbuddy/internal/permify"
	pgeventing "github.com/rsmidt/soccerbuddy/internal/postgres/eventing"
	"github.com/rsmidt/soccerbuddy/internal/projector"
	rdeventing "github.com/rsmidt/soccerbuddy/internal/redis/eventing"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	grpc2 "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

// https://docs.permify.co/api-reference/tenancy/create-tenant
const permifyTenantID = "t1"

func main() {
	ctx := context.Background()
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	conf, err := getConf()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := run(ctx, conf, slog.New(handler)); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	os.Exit(1)
}

func run(ctx context.Context, c *config.Config, log *slog.Logger) (err error) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	shutdown, err := tracing.SetupOtelSDK(ctx)
	if err != nil {
		return err
	}
	defer func() {
		// TODO: Perform shutdown with timeout.
		err = errors.Join(err, shutdown(context.Background()))
	}()

	// Setup the postgres connection.
	dbconfig, err := pgxpool.ParseConfig(c.EventJournal.PG.ConnStr())
	if err != nil {
		return err
	}
	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, "SET application_name = 'soccerbuddy'")
		if err != nil {
			return err
		}
		pgxdecimal.Register(conn.TypeMap())
		return err
	}
	pool, err := pgxpool.NewWithConfig(ctx, dbconfig)
	if err != nil {
		return fmt.Errorf("failed to create pg pool: %w", err)
	}
	err = pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
		return runMigration(ctx, conn.Conn())
	})
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Setup redis.
	rdOpts := rueidis.ClientOption{
		InitAddress: []string{c.EventJournal.Redis.Host},
	}
	rdClient, err := rueidisotel.NewClient(rdOpts)
	if err != nil {
		return fmt.Errorf("failed to create redis client: %w", err)
	}
	rdLocker, err := rueidislock.NewLocker(rueidislock.LockerOption{
		ClientOption:   rdOpts,
		KeyMajority:    1,
		NoLoopTracking: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create redis locker: %w", err)
	}

	// Setup permify authorizer.
	client, err := setupPermifyClient(ctx, c)
	if err != nil {
		return fmt.Errorf("failed to setup permify client: %w", err)
	}
	authorizer := permify.NewAuthorizer(log, client, permifyTenantID)
	relationStore := permify.NewRelationStore(log, client)

	// Setup event store.
	eventCrypto := pgeventing.NewEventCrypto(pool)
	es := pgeventing.NewEventStore(log, pool, eventregistry.Default, eventCrypto, relationStore)

	// Setup application.
	repos := assembleRepositories(es)
	cmds := commands.NewCommands(log, es, authorizer, rdClient, repos)
	qs := queries.NewQueries(log, es, authorizer, rdClient, repos)

	// Setup projectors.
	ps := pgeventing.NewProjectorSupervisor(log, pool, es)
	rds := rdeventing.NewProjectorSupervisor(log, es, rdClient, rdLocker)
	supervisors := projector.Supervisors{Postgres: ps, Redis: rds}
	if err := supervisors.Register(ctx, relationStore, rdClient); err != nil {
		return fmt.Errorf("failed to register and init projectors: %v", err)
	}
	supervisors.Enable()

	pgEn := pgeventing.NewEventNotifier(log, pool)
	pgEn.AddListener(ps)
	pgEn.AddListener(rds)
	go func() {
		// Delay starting of the event listener to allow for initial triggers to not compete directly.
		// This is not strictly necessary, but will ease the startup a bit.
		time.Sleep(1 * time.Second)
		if err := pgEn.Start(ctx); err != nil {
			log.Error("Event listener failed", "error", err)
		}
	}()

	// Start the projection polling loop.
	go func() {
		interval := c.Projection.PollingInterval
		timer := time.NewTimer(interval)
		defer timer.Stop()

		log.Info("Starting projection polling.")

		for {
			select {
			case <-ctx.Done():
				log.Debug("Stopping projection polling loop.")
				return
			case <-timer.C:
				log.Debug("Triggering projection polling run.")
				var wg sync.WaitGroup
				wg.Add(2)
				go func() {
					defer wg.Done()
					rds.Trigger(ctx)
				}()
				go func() {
					defer wg.Done()
					ps.Trigger(ctx)
				}()
				wg.Wait()

				timer.Reset(interval)
			}
		}
	}()

	// On every persisted event, we want to trigger the account permission projector.
	// Here, we can't really deal with eventual consistency.
	// The projector will wait for any locks on the projection to avoid race conditions when any
	// NOTIFY triggered projector is faster.
	es.AddHook(eventing.NewPostPersistHook(func(ctx context.Context) error {
		ps.Trigger(ctx, projector.PermissionProjectorName)
		return nil
	}))

	// Create root account if it doesn't exist.
	if err := cmds.CreateRootAccount(ctx, commands.CreateRootAccountCommand{
		Email:     c.Setup.Root.Email,
		Password:  c.Setup.Root.Password,
		FirstName: c.Setup.Root.FirstName,
		LastName:  c.Setup.Root.LastName,
	}); err != nil && !errors.Is(err, domain.ErrRootAccountAlreadyInitialized) {
		return err
	}

	// Setup the main http server including grpc via connect.
	grpcServer := grpc.NewServer(cmds, qs, log)
	mux := http.NewServeMux()
	if err := grpcServer.Register(mux); err != nil {
		return err
	}
	rootHandler := h2c.NewHandler(mux, &http2.Server{})

	srv := &http.Server{
		Addr:        c.Host,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
		Handler:     rootHandler,
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	log.Info("Starting application")

	select {
	case <-ctx.Done():
		stop()
	case err = <-srvErr:
		return
	}

	// TODO: Perform shutdown with timeout.
	err = srv.Shutdown(context.Background())
	return
}

func runMigration(ctx context.Context, conn *pgx.Conn) error {
	migrator, err := migrate.NewMigrator(ctx, conn, "public.schema_version")
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	// Load migration directory as a fs.FS.
	migrationFS := os.DirFS("migrations")
	if err := migrator.LoadMigrations(migrationFS); err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}
	return migrator.Migrate(ctx)
}

func setupPermifyClient(ctx context.Context, c *config.Config) (*permify_grpc.Client, error) {
	client, err := permify_grpc.NewClient(
		permify_grpc.Config{
			Endpoint: c.Permify.Host,
		},
		grpc2.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func getConf() (*config.Config, error) {
	viper.SetEnvPrefix("soccerbuddy")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	viper.SetConfigName(".local.dev")
	if err := viper.MergeInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("failed to find local dev conf: %w", err)
		}
	}
	viper.SetConfigName("config")
	if err := viper.MergeInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("failed to find conf: %w", err)
		}
	}
	return config.NewConfig(viper.GetViper())
}

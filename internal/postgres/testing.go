package postgres

import (
	"context"
	"fmt"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"github.com/rsmidt/soccerbuddy/internal/core"
	"github.com/rsmidt/soccerbuddy/internal/core/idgen"
	"log"
	"log/slog"
	"os"
	"path"
	"strings"
	"sync"
)

type DbCleanup func()

type dbPoolGetter interface {
	getTestPool() (*pgxpool.Pool, DbCleanup)
}

type pooledDbGetter struct {
	pool *pgxpool.Pool
}

func (p *pooledDbGetter) getTestPool() (*pgxpool.Pool, DbCleanup) {
	dbname := fmt.Sprintf("test_%s", strings.ToLower(idgen.NewString()))
	core.Must(p.pool.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
		_, err := conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s TEMPLATE soccerbuddy", dbname))
		return err
	}))
	config := core.Must2(pgxpool.ParseConfig(fmt.Sprintf("postgres://soccerbuddy:soccerbuddy@localhost:34344/%s", dbname)))
	config.MaxConns = 5
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, "SET application_name = 'soccerbuddy'")
		if err != nil {
			return err
		}

		pgxdecimal.Register(conn.TypeMap())

		return err
	}
	pool := core.Must2(pgxpool.NewWithConfig(context.Background(), config))
	return pool, func() {
		pool.Close()
		log.Println("Dropping test db")
		core.Must(p.pool.AcquireFunc(context.Background(), func(conn *pgxpool.Conn) error {
			_, err := conn.Exec(context.Background(), fmt.Sprintf("DROP DATABASE %s",
				dbname))
			return err
		}))
	}
}

var (
	dbGetterInit sync.Once
	dbGetter     dbPoolGetter
)

func GetTestPool() (*pgxpool.Pool, DbCleanup) {
	ctx := context.Background()

	dbGetterInit.Do(func() {
		config := core.Must2(pgxpool.ParseConfig("postgres://soccerbuddy:soccerbuddy@localhost:34344/soccerbuddy"))
		config.MaxConns = 1
		pool := core.Must2(pgxpool.NewWithConfig(ctx, config))
		core.Must(pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
			migrator := core.Must2(migrate.NewMigrator(ctx, conn.Conn(), "public.schema_version"))
			migrationFS := os.DirFS(path.Join(core.Root, "migrations"))
			core.Must(migrator.LoadMigrations(migrationFS))
			core.Must(migrator.Migrate(ctx))

			core.Must2(conn.Exec(ctx, "ALTER DATABASE soccerbuddy is_template = true"))
			return nil
		}))
		// Drop all existing test databases.
		rows := core.Must2(pool.Query(ctx, "SELECT datname FROM pg_database WHERE datname LIKE 'test_%'"))
		dbnames := core.Must2(pgx.CollectRows(rows, func(row pgx.CollectableRow) (string, error) {
			var dbname string
			err := row.Scan(&dbname)
			return dbname, err
		}))
		for _, dbname := range dbnames {
			core.Must(pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
				_, err := conn.Exec(ctx, fmt.Sprintf("DROP DATABASE %s", dbname))
				return err
			}))
		}
		dbGetter = &pooledDbGetter{pool: pool}
	})

	return dbGetter.getTestPool()
}

func GetTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

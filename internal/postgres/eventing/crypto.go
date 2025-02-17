package eventing

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rsmidt/soccerbuddy/internal/eventing"
	"github.com/rsmidt/soccerbuddy/internal/postgres"
	"github.com/rsmidt/soccerbuddy/internal/tracing"
	"io"
	"maps"
	"slices"
)

type (
	keysByOwner map[eventing.AggregateID][]byte
)

type pgEventCrypto struct {
	pool *pgxpool.Pool
}

func NewEventCrypto(pool *pgxpool.Pool) eventing.EventCrypto {
	return &pgEventCrypto{
		pool: pool,
	}
}

func (p *pgEventCrypto) EncryptEvents(ctx context.Context, events []eventing.Event) error {
	ctx, span := tracing.Tracer.Start(ctx, "pgEventCrypto.EncryptEvents")
	defer span.End()

	owners, encEvents, err := p.filterEncryptedEvents(events)
	if err != nil {
		return fmt.Errorf("failed to find encrypted values: %w", err)
	}
	if len(encEvents) == 0 {
		return nil
	}
	keys, err := p.loadOrCreateKeys(ctx, owners)
	if err != nil {
		return fmt.Errorf("failed to load key: %w", err)
	}
	// Encrypt the values.
	transformer := &aesEncryptor{keys: keys}
	for _, event := range encEvents {
		if err := event.AcceptCrypto(transformer); err != nil {
			return fmt.Errorf("failed to encrypt value: %w", err)
		}
	}
	return nil
}

func (p *pgEventCrypto) DecryptEvents(ctx context.Context, events []eventing.Event) error {
	ctx, span := tracing.Tracer.Start(ctx, "pgEventCrypto.DecryptEvents")
	defer span.End()

	owners, encEvents, err := p.filterEncryptedEvents(events)
	if err != nil {
		return fmt.Errorf("failed to find encrypted values: %w", err)
	}
	if len(encEvents) == 0 {
		return nil
	}
	key, err := p.maybeLoadKeys(ctx, owners)
	if err != nil {
		return fmt.Errorf("failed to load key: %w", err)
	}
	// Decrypt the values.
	transformer := &aesDecryptor{keys: key}
	for _, event := range encEvents {
		if err := event.AcceptCrypto(transformer); err != nil {
			return fmt.Errorf("failed to decrypt value: %w", err)
		}
	}
	return nil
}

func (p *pgEventCrypto) filterEncryptedEvents(events []eventing.Event) ([]eventing.AggregateID, []eventing.EncryptedEvent, error) {
	var encryptedEvents []eventing.EncryptedEvent
	owners := make(map[eventing.AggregateID]struct{})
	for _, event := range events {
		event, ok := event.(eventing.EncryptedEvent)
		if !ok {
			continue
		}
		for _, owner := range event.DeclareOwners() {
			owners[owner] = struct{}{}
		}
		encryptedEvents = append(encryptedEvents, event)
	}

	return slices.Collect(maps.Keys(owners)), encryptedEvents, nil
}

func (p *pgEventCrypto) loadOrCreateKeys(ctx context.Context, owners []eventing.AggregateID) (keysByOwner, error) {
	db := postgres.GetDBFromContext(ctx, p.pool)
	rows, err := db.Query(ctx, "SELECT owner_id, key FROM keys WHERE owner_id = ANY ($1)", owners)
	if err != nil {
		return nil, err
	}
	var (
		key        []byte
		ownerID    eventing.AggregateID
		keyByOwner = make(keysByOwner)
	)
	_, err = pgx.ForEachRow(rows, []any{&ownerID, &key}, func() error {
		keyByOwner[ownerID] = key
		return nil
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if len(keyByOwner) == len(owners) {
		return keyByOwner, nil
	}
	missingKeysToInsert := make(map[eventing.AggregateID][]byte)
	for _, owner := range owners {
		if _, ok := keyByOwner[owner]; !ok {
			key := make([]byte, 32)
			if _, err := io.ReadFull(rand.Reader, key); err != nil {
				return nil, err
			}
			missingKeysToInsert[owner] = key
		}
	}
	if len(missingKeysToInsert) > 0 {
		err = pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
			_, err = tx.CopyFrom(ctx, pgx.Identifier{"keys"}, []string{"owner_id", "key"}, pgx.CopyFromSlice(len(missingKeysToInsert), func(i int) ([]interface{}, error) {
				owner := owners[i]
				key := missingKeysToInsert[owner]
				return []interface{}{owner, key}, nil
			}))
			if err != nil {
				return fmt.Errorf("failed to insert missing keys: %w", err)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	for owner, key := range missingKeysToInsert {
		keyByOwner[owner] = key
	}
	return keyByOwner, nil
}

func (p *pgEventCrypto) maybeLoadKeys(ctx context.Context, owners []eventing.AggregateID) (keysByOwner, error) {
	db := postgres.GetDBFromContext(ctx, p.pool)
	rows, err := db.Query(ctx, "SELECT owner_id, key FROM keys WHERE owner_id = ANY ($1)", owners)
	if err != nil {
		return nil, err
	}
	var (
		key        []byte
		ownerID    eventing.AggregateID
		keyByOwner = make(keysByOwner)
	)
	_, err = pgx.ForEachRow(rows, []any{&ownerID, &key}, func() error {
		keyByOwner[ownerID] = key
		return nil
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return keyByOwner, nil
}

type aesEncryptor struct {
	keys keysByOwner
}

func (a *aesEncryptor) Transform(owner eventing.AggregateID, value *eventing.EncryptedString) error {
	key, ok := a.keys[owner]
	if !ok {
		return fmt.Errorf("key for owner should have been generated already: %s", owner)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	ciphertext := aead.Seal(nil, nonce, []byte(value.Value), nil)
	value.Value = base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))
	return nil
}

func (a *aesEncryptor) TransformWithDefault(owner eventing.AggregateID, value *eventing.EncryptedString, defaultValue string) error {
	return a.Transform(owner, value)
}

type aesDecryptor struct {
	keys keysByOwner
}

func (a *aesDecryptor) Transform(owner eventing.AggregateID, value *eventing.EncryptedString) error {
	return a.TransformWithDefault(owner, value, "")
}

func (a *aesDecryptor) TransformWithDefault(owner eventing.AggregateID, value *eventing.EncryptedString, defaultValue string) error {
	key, ok := a.keys[owner]
	if !ok {
		value.IsShredded = true
		value.Value = defaultValue
		return nil
	}
	ciphertext, err := base64.StdEncoding.DecodeString(value.Value)
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonceSize := aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}
	value.Value = string(plaintext)
	return nil
}

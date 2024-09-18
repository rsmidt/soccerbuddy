package eventing

import (
	"context"
	"errors"
)

var (
	ErrNoDecryptionKey           = errors.New("no decryption key found")
	ErrAggregateIsCryptoShredded = errors.New("aggregate is crypto shredded")
)

type EncryptedString struct {
	Value      string `json:"value"`
	IsShredded bool   `json:"-"`
}

func NewEncryptedString(value string) EncryptedString {
	return EncryptedString{
		Value:      value,
		IsShredded: false,
	}
}

type CryptoTransformer interface {
	Transform(owner AggregateID, value *EncryptedString) error
	TransformWithDefault(owner AggregateID, value *EncryptedString, defaultValue string) error
}

type EncryptedEvent interface {
	DeclareOwners() []AggregateID
	AcceptCrypto(transformer CryptoTransformer) error
}

type EventEncryptor interface {
	// EncryptEvents encrypts the values in the event in place.
	EncryptEvents(ctx context.Context, events []Event) error
}

type EventDecrypter interface {
	// DecryptEvents decrypts the encrypted values in the event in place.
	DecryptEvents(ctx context.Context, events []Event) error
}

type EventCrypto interface {
	EventEncryptor
	EventDecrypter
}

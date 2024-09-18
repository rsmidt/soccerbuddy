package idgen

import "github.com/segmentio/ksuid"

type KSUIDGenerator struct {
}

func (k *KSUIDGenerator) Gen() string {
	return ksuid.New().String()
}

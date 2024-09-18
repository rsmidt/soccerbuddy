package redis

import (
	"encoding/json"
	"github.com/redis/rueidis"
	"strings"
)

// UnmarshalDocs unmarshals a slice of rueidis.FtSearchDoc into a slice of T.
func UnmarshalDocs[T any](docs []rueidis.FtSearchDoc) ([]*T, error) {
	res := make([]*T, len(docs))
	for i, doc := range docs {
		var t T
		raw := doc.Doc["$"]
		if err := json.NewDecoder(strings.NewReader(raw)).Decode(&t); err != nil {
			return nil, err
		}
		res[i] = &t
	}
	return res, nil
}

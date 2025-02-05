package redis

import (
	"encoding/json"
	"github.com/redis/rueidis"
	"strings"
)

// UnmarshalDocs unmarshals a slice of rueidis.FtSearchDoc into a slice of T.
func UnmarshalDocs[T any](docs []rueidis.FtSearchDoc) ([]*T, error) {
	res := make([]*T, 0, len(docs))
	for _, doc := range docs {
		var t []*T
		raw := doc.Doc["$"]
		if err := json.NewDecoder(strings.NewReader(raw)).Decode(&t); err != nil {
			return nil, err
		}
		if t[0] != nil {
			res = append(res, t[0])
		}
	}
	return res, nil
}

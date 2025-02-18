package redis

import (
	"encoding/json"
	"github.com/redis/rueidis"
	"strings"
	"time"
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

// FlattenToString returns:
// - an empty string if s is empty,
// - the first element if s is a JSON array (e.g., '["first", "second"]'),
// - or s itself if it's not an array.
func FlattenToString(s string) string {
	if s == "" {
		return ""
	}

	// Check if the string might be a JSON array.
	if strings.HasPrefix(s, "[") {
		var arr []string
		if err := json.Unmarshal([]byte(s), &arr); err == nil {
			if len(arr) > 0 {
				return arr[0]
			}
			return ""
		}
		// If unmarshaling fails, we assume s is not a valid JSON array.
	}
	return s
}

// FlattenToTime parses the time.
func FlattenToTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, FlattenToString(s))
	if err != nil {
		return time.Time{}
	}
	return t
}

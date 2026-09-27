package helper

import (
	"strings"

	"github.com/google/uuid"
)

func ParseUUIDs(ids []string) ([]uuid.UUID, error) {
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		u, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func HasNullByte(ss ...string) bool {
	for _, s := range ss {
		if strings.IndexByte(s, 0) >= 0 {
			return true
		}
	}
	return false
}
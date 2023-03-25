package creators

import (
	"github.com/google/uuid"
)

func UuidString() string {
	id, _ := uuid.NewRandom()
	return id.String()
}

package base

import (
	"time"

	"github.com/google/uuid"
)

type TBase struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdat"`
}

type TBaseWithName struct {
	TBase
	Name string `json:"name"`
}

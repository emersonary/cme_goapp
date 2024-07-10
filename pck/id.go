package pck

import "github.com/google/uuid"

type UUID = uuid.UUID

func NewID() UUID {
	return UUID(uuid.New())
}

func ParseID(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	return UUID(id), err
}

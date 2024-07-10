package handler

import (
	"github.com/emersonary/go-authentication/model/base"
	"github.com/emersonary/go-authentication/pck"
)

type HandlerDBUUID = interface {
	FindById(id pck.UUID) (*any, error)
	DeleteById(id pck.UUID) error
}

type HandlerDBInterface = interface {
	Insert(entity any) (base.TBase, error)
	Update(entity any)
	Store(entity any)
	List() []any
}

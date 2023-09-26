package repository

import (
	"context"

	"github.com/otakakot/gosql/internal/domain/model"
)

type Transactor interface {
	Transact(context.Context, ...string) error
}

type UserQuery interface {
	Find(context.Context, string) (*model.User, error)
}

type UserCommand interface {
	Create(context.Context, model.User) ([]string, error)
}

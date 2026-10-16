package input

import (
	"context"
)

type BaseManagementService[T any] interface {
	GetAll(ctx context.Context) (*[]T, error)
	GetById(ctx context.Context, id int32) (*T, error)
}

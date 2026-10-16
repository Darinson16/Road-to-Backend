package output

import (
	"context"
)

type BaseManagementRepository[T any] interface {
	GetAll(ctx context.Context) (*[]T, error)
	GetById(ctx context.Context, id int32) (*T, error)
}

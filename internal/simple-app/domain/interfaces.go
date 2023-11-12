// Package domain contains application domain layer
package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

//go:generate mockery --name=ApplicationInterface
type ApplicationInterface interface {
	GetUser(ctx context.Context, id uuid.UUID) (name string, err error)
	CreateUser(ctx context.Context, name string) (id uuid.UUID, err error)
	UpdateUser(ctx context.Context, id uuid.UUID, name string) (err error)
	DeleteUser(ctx context.Context, id uuid.UUID) (err error)
}

var ErrorNotFound = fmt.Errorf("not found")

//go:generate mockery --name=UserStorage
type UserStorage interface {
	Store(ctx context.Context, id, name string) (err error)
	Read(ctx context.Context, id string) (name string, err error)
	Delete(ctx context.Context, id string) (err error)
}

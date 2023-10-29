// Package application contains application layer
package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/adlandh/acorn-simple-app/internal/simple-app/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var _ domain.ApplicationInterface = (*Application)(nil)

type Application struct {
	logger  *zap.Logger
	storage domain.UserStorage
}

func NewApplication(logger *zap.Logger, storage domain.UserStorage) *Application {
	return &Application{
		logger:  logger,
		storage: storage,
	}
}

func (a Application) GetUser(ctx context.Context, id uuid.UUID) (name string, err error) {
	strID := id.String()

	name, err = a.storage.Read(ctx, strID)
	if err == nil || errors.Is(err, domain.ErrorNotFound) {
		return
	}

	a.logger.Error("error getting user", zap.String("id", strID), zap.Error(err))

	return name, fmt.Errorf("error getting message: %s", id)
}

func (a Application) CreateUser(ctx context.Context, name string) (id uuid.UUID, err error) {
	id, err = uuid.NewUUID()
	if err != nil {
		a.logger.Error("error generating uuid", zap.Error(err), zap.String("name", name))
		return id, fmt.Errorf("error generating id")
	}

	err = a.storage.Store(ctx, id.String(), name)
	if err != nil {
		a.logger.Error("error creating user", zap.Error(err), zap.String("id", id.String()), zap.String("name", name))
		return id, fmt.Errorf("error creating user")
	}

	return
}

func (a Application) UpdateUser(ctx context.Context, id uuid.UUID, name string) (err error) {
	strID := id.String()

	_, err = a.storage.Read(ctx, strID)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			return
		}

		a.logger.Error("error updating user", zap.Error(err), zap.String("id", strID), zap.String("name", name))

		return fmt.Errorf("error getting user")
	}

	err = a.storage.Store(ctx, strID, name)
	if err != nil {
		a.logger.Error("error updating user", zap.Error(err), zap.String("id", strID), zap.String("name", name))

		return fmt.Errorf("error updating user")
	}

	return
}

func (a Application) DeleteUser(ctx context.Context, id uuid.UUID) (err error) {
	strID := id.String()

	err = a.storage.Delete(ctx, strID)
	if err != nil {
		if errors.Is(err, domain.ErrorNotFound) {
			return
		}

		a.logger.Error("error deleting user", zap.Error(err), zap.String("id", strID))

		return fmt.Errorf("error deleting user")
	}

	return
}

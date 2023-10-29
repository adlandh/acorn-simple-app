package application

import (
	"context"
	"testing"

	"github.com/adlandh/acorn-simple-app/internal/simple-app/domain/mocks"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestCreateGetUpdateAndDeleteUser(t *testing.T) {
	storage := new(mocks.UserStorage)
	logger := zaptest.NewLogger(t)
	app := NewApplication(logger, storage)
	ctx := context.Background()

	t.Run("create user", func(t *testing.T) {
		name := gofakeit.Username()
		storage.On("Store", ctx, mock.Anything, name).Return(nil).Once()
		id, err := app.CreateUser(ctx, name)
		require.NoError(t, err)
		require.NotEmpty(t, id)
	})

	t.Run("get user", func(t *testing.T) {
		id, err := uuid.NewUUID()
		require.NoError(t, err)
		name := gofakeit.Username()
		storage.On("Read", ctx, id.String()).Return(name, nil).Once()
		storedName, err := app.GetUser(ctx, id)
		require.NoError(t, err)
		require.Equal(t, name, storedName)
	})

	t.Run("update user", func(t *testing.T) {
		id, err := uuid.NewUUID()
		require.NoError(t, err)
		name := gofakeit.Username()
		storage.On("Read", ctx, id.String()).Return(name, nil).Once()
		newName := gofakeit.Username()
		storage.On("Store", ctx, id.String(), newName).Return(nil).Once()

		err = app.UpdateUser(ctx, id, newName)
		require.NoError(t, err)
	})

	t.Run("delete user", func(t *testing.T) {
		id, err := uuid.NewUUID()
		require.NoError(t, err)
		storage.On("Delete", ctx, id.String()).Return(nil).Once()
		err = app.DeleteUser(ctx, id)
		require.NoError(t, err)
	})

	storage.AssertExpectations(t)
}

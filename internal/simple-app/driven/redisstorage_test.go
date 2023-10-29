package driven

import (
	"context"
	"testing"
	"time"

	"github.com/adlandh/acorn-simple-app/internal/simple-app/config"
	"github.com/adlandh/acorn-simple-app/internal/simple-app/domain"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx/fxtest"
)

type RedisStorageTestSuite struct {
	suite.Suite
	storage *RedisStorage
	id      string
	name    string
}

func (s *RedisStorageTestSuite) SetupSuite() {
	ctx := context.Background()

	s.id = gofakeit.UUID()
	s.name = gofakeit.Username()

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Ready to accept connections").WithStartupTimeout(3*time.Minute),
			wait.ForListeningPort("6379/tcp").WithStartupTimeout(3*time.Minute),
		),
		Name: "redis",
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err)

	redisPort, err := container.MappedPort(ctx, "6379")
	s.Require().NoError(err)

	port := redisPort.Port()

	host, err := container.Host(ctx)
	s.Require().NoError(err)

	if host == "" {
		host = "localhost"
	}

	lc := fxtest.NewLifecycle(s.T())

	s.storage, err = NewRedisStorage(
		lc,
		&config.Config{
			Redis: config.RedisConfig{
				URL:    "redis://" + host + ":" + port,
				Prefix: gofakeit.Word(),
			},
		})
	s.Require().NoError(err)

	err = lc.Start(ctx)
	s.Require().NoError(err)
}

func (s *RedisStorageTestSuite) Test1Store() {
	err := s.storage.Store(context.Background(), s.id, s.name)
	s.Require().NoError(err)
}

func (s *RedisStorageTestSuite) Test2Read() {
	name, err := s.storage.Read(context.Background(), s.id)
	s.Require().NoError(err)
	s.Require().Equal(s.name, name)
}

func (s *RedisStorageTestSuite) Test3NotFound() {
	id := gofakeit.UUID()

	_, err := s.storage.Read(context.Background(), id)
	s.Require().Error(err)
	s.Require().Equal(domain.ErrorNotFound, err)
}

func (s *RedisStorageTestSuite) Test4Delete() {
	err := s.storage.Delete(context.Background(), s.id)
	s.Require().NoError(err)

	_, err = s.storage.Read(context.Background(), s.id)
	s.Require().Error(err)
	s.Require().Equal(domain.ErrorNotFound, err)
}

func TestRedisStorage(t *testing.T) {
	suite.Run(t, new(RedisStorageTestSuite))
}

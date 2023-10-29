// Package driven contains driven layer
package driven

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/adlandh/acorn-simple-app/internal/simple-app/config"
	"github.com/adlandh/acorn-simple-app/internal/simple-app/domain"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var _ domain.UserStorage = (*RedisStorage)(nil)

type RedisStorage struct {
	client *redis.Client
	prefix string
}

func NewRedisStorage(lc fx.Lifecycle, cfg *config.Config) (*RedisStorage, error) {
	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return nil, fmt.Errorf("error parsing redis url: %w", err)
	}

	r := &RedisStorage{
		client: redis.NewClient(opt),
		prefix: cfg.Redis.Prefix,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err := r.client.Set(ctx, r.genID("ping"), "pong", 1*time.Millisecond).Err()
			if err != nil {
				return fmt.Errorf("error connecting to redis: %w", err)
			}
			return nil
		},
	})

	return r, nil
}

func (r RedisStorage) Store(ctx context.Context, id, name string) (err error) {
	err = r.client.Set(ctx, r.genID(id), name, 0).Err()
	if err != nil {
		err = fmt.Errorf("error storing to redis: %w", err)
	}

	return
}

func (r RedisStorage) Read(ctx context.Context, id string) (name string, err error) {
	id = r.genID(id)

	name, err = r.client.Get(ctx, id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = domain.ErrorNotFound

			return
		}

		err = fmt.Errorf("error reading from redis: %w", err)

		return
	}

	return
}

func (r RedisStorage) Delete(ctx context.Context, id string) (err error) {
	err = r.client.Del(ctx, r.genID(id)).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err = domain.ErrorNotFound

			return
		}

		err = fmt.Errorf("error delete from redis: %w", err)

		return
	}

	return
}

func (r RedisStorage) genID(id string) string {
	return r.prefix + "::" + id
}

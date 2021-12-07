package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/kmhebb/serverExample/internal/log"
)

type EmailValidationCache struct {
	pool   *redis.Pool
	logger log.Logger
}

func NewEmailValidationCache(
	pool *redis.Pool,
	logger log.Logger,
) EmailValidationCache {
	return EmailValidationCache{
		pool:   pool,
		logger: logger,
	}
}

func (c EmailValidationCache) DeleteValidationCode(ctx context.Context, email string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", email)
	return err
}

func (c EmailValidationCache) GetValidationCode(ctx context.Context, email string) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	code, err := redis.String(conn.Do("GET", email))
	if err != nil {
		if err == redis.ErrNil {
			return "", fmt.Errorf("no validation code for email")
		}
		return "", err
	}
	return code, nil
}

func (c EmailValidationCache) SaveValidationCode(ctx context.Context, email, code string) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("SETEX", email, time.Hour.Seconds(), code))
	return err
}

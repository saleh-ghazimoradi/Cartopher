package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Options func(*Redis)

func WithHost(host string) Options {
	return func(r *Redis) {
		r.Host = host
	}
}

func WithPort(port string) Options {
	return func(r *Redis) {
		r.Port = port
	}
}

func WithPassword(password string) Options {
	return func(r *Redis) {
		r.Password = password
	}
}

func WithDB(db int) Options {
	return func(r *Redis) {
		r.DB = db
	}
}

func (r *Redis) uri() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

func (r *Redis) Connect(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     r.uri(),
		Password: r.Password,
		DB:       r.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedis(opts ...Options) *Redis {
	r := &Redis{}
	for _, o := range opts {
		o(r)
	}
	return r
}

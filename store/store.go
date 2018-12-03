package store

import (
	"context"
	"errors"
)

var (
	ErrStoreNotFound       = errors.New("not found")
	ErrStoreNotStringValue = errors.New("value is not string")
)

type Provider interface {
	Store(namespace string) Store
}

type Store interface {
	Get(ctx context.Context, key string) (interface{}, error)
	GetString(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val interface{}) error
}

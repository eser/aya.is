package caching

import (
	"context"
	"errors"
	"fmt"

	"github.com/eser/aya.is/services/pkg/lib/vars"
)

var (
	ErrCannotGetFromCache     = errors.New("cannot get from cache")
	ErrCannotSetToCache       = errors.New("cannot set to cache")
	ErrCannotExecuteCachingFn = errors.New("cannot execute caching function")
)

type Cache struct {
	getter func(ctx context.Context, key string, target any) (bool, error)
	setter func(ctx context.Context, key string, value any) error
}

func NewCache(
	getter func(ctx context.Context, key string, target any) (bool, error),
	setter func(ctx context.Context, key string, value any) error,
) *Cache {
	return &Cache{
		getter: getter,
		setter: setter,
	}
}

func (c *Cache) Get(ctx context.Context, target any, key string) (bool, error) {
	isSet, err := c.getter(ctx, key, target)
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrCannotGetFromCache, err)
	}

	if !isSet {
		// fmt.Println("cache miss ", key) //nolint:forbidigo
		return false, nil
	}

	return true, nil
}

func (c *Cache) Set(ctx context.Context, key string, value any) error {
	err := c.setter(ctx, key, value)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotSetToCache, err)
	}

	return nil
}

func (c *Cache) Execute(
	ctx context.Context,
	key string,
	target any,
	fn func(ctx context.Context) (any, error), //nolint:varnamelen
) error {
	isGot, err := c.Get(ctx, target, key)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotExecuteCachingFn, err)
	}

	if isGot {
		return nil
	}

	value, err := fn(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotExecuteCachingFn, err)
	}

	err = vars.SetValue(target, value)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotExecuteCachingFn, err)
	}

	err = c.Set(ctx, key, value)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCannotExecuteCachingFn, err)
	}

	return nil
}

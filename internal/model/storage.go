package model

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type Storage[T any] interface {
	Read(ctx context.Context, id uuid.UUID) (T, error)
	Create(ctx context.Context, t T) (uuid.UUID, error)
	Update(ctx context.Context, t T) (T, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]T, error)
	First(ctx context.Context, f func(T) bool) (T, error)
}

var (
	ErrNotFound = errors.New("Not Found")
	ErrConflict = errors.New("Conflict")
)

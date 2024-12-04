package storage

import (
	"github.com/google/uuid"
)

type Storage[T any] interface {
	Get(id uuid.UUID) (T, error)
	Create(model *T) error
	Delete(model *T) error
}

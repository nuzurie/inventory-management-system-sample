package domain

import (
	"context"
	"time"
)

type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ItemUseCase interface {
	GetAll(ctx context.Context, count int, offset int, filter Specification) ([]Item, error)
	GetOne(ctx context.Context, id string) (*Item, error)
	Create(ctx context.Context, item *Item) (*Item, error)
	Update(ctx context.Context, item *Item) (*Item, error)
	Delete(ctx context.Context, id string) error
}

type ItemRepository interface {
	GetAll(ctx context.Context, count int, offset int, filter Specification) ([]Item, error)
	GetOne(ctx context.Context, id string) (*Item, error)
	Save(ctx context.Context, item *Item) (*Item, error)
	Edit(ctx context.Context, item *Item) (*Item, error)
	Delete(ctx context.Context, id string) error
}

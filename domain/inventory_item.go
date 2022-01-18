package domain

import (
	"context"
	"time"
)

type InventoryItem struct {
	ID        string    `json:"id"`
	Item      Item      `json:"item"`
	Quantity  int       `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at"`
}

type InventoryUseCase interface {
	GetInventoryForItem(ctx context.Context, itemID string) (*InventoryItem, error)
	GetAll(ctx context.Context, count int, offset int, filter Specification) ([]InventoryItem, error)
	UpdateInventoryItem(ctx context.Context, item *InventoryItem) (*InventoryItem, error)
	DeleteItem(ctx context.Context, item *InventoryItem) error
}

type InventoryRepository interface {
	GetInventoryForItem(ctx context.Context, itemID string) (*InventoryItem, error)
	GetAll(ctx context.Context, count int, offset int, filter Specification) ([]InventoryItem, error)
	Save(ctx context.Context, item *InventoryItem) (*InventoryItem, error)
	Edit(ctx context.Context, item *InventoryItem) (*InventoryItem, error)
	DeleteItem(ctx context.Context, item *InventoryItem) error
}

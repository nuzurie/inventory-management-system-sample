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
	// GetInventoryForItem to test if an item has any stock in the inventory
	GetInventoryForItem(ctx context.Context, itemID string) (*InventoryItem, error)
	GetAll(ctx context.Context, count int, offset int, filter InventorySpecification) ([]InventoryItem, error)
	UpdateInventoryItem(ctx context.Context, item *InventoryItem) (*InventoryItem, error)
	DeleteItem(ctx context.Context, id string) error
}

type InventoryRepository interface {
	GetInventoryForItem(ctx context.Context, itemID string) (*InventoryItem, error)
	GetAll(ctx context.Context, count int, offset int, filter InventorySpecification) ([]InventoryItem, error)
	GetByID(ctx context.Context, id string) (*InventoryItem, error)
	Save(ctx context.Context, item *InventoryItem) (*InventoryItem, error)
	Edit(ctx context.Context, item *InventoryItem) (*InventoryItem, error)
	DeleteItem(ctx context.Context, id string) error
}

package usecase

import (
	"context"
	"fmt"
	"github.com/nuzurie/shopify/domain"
	"github.com/nuzurie/shopify/utils/errors"
	"github.com/google/uuid"
	"log"
	"reflect"
	"time"
)

type itemUseCase struct {
	itemRepository 	domain.ItemRepository
	timeout 		time.Duration
}

func NewItemUseCase(repository domain.ItemRepository, timeout time.Duration) domain.ItemUseCase {
	return &itemUseCase{itemRepository: repository, timeout: timeout}
}

func (i *itemUseCase) GetAll(ctx context.Context, count int, offset int, filter domain.Specification) ([]domain.Item, error) {
	c, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	items, err := i.itemRepository.GetAll(c, count, offset, filter)
	if err != nil {
		return nil, err
	}

	return items, err
}

func (i *itemUseCase) GetOne(ctx context.Context, id string) (*domain.Item, error) {
	c, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	item, err := i.itemRepository.GetOne(c, id)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(item, &domain.Item{}) {
		return nil, errors.NewBadRequestError("no such item exists")
	}

	return item, nil
}

func (i *itemUseCase) Create(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	c, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	item.ID = uuid.NewString()
	createdItem, err := i.itemRepository.Save(c, item)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to create item %s at %s", item.ID, time.Now()))
		return nil, err
	}

	return createdItem, nil
}

func (i *itemUseCase) Update(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	c, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	existingItem, err := i.itemRepository.GetOne(c, item.ID)
	if err != nil {
		return nil, err
	}
	if reflect.DeepEqual(existingItem, &domain.Item{}) {
		return nil, errors.NewBadRequestError("no such item exists")
	}

	item.UpdatedAt = time.Now()
	var updated *domain.Item
	updated, err = i.itemRepository.Edit(c, item)
	if err != nil {
		return nil, err
	}

	return updated, err
}

func (i itemUseCase) Delete(ctx context.Context, id string) error {
	c, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	existingItem, err := i.itemRepository.GetOne(c, id)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(existingItem, &domain.Item{}) {
		return errors.NewBadRequestError("no such item exists")
	}

	return nil
}
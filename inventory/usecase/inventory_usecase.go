package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/nuzurie/shopify/domain"
	"github.com/nuzurie/shopify/utils/errors"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

type inventoryUseCase struct {
	itemRepository      domain.ItemRepository
	inventoryRepository domain.InventoryRepository
	timeout             time.Duration
}

func NewInventoryUseCase(itemRepository domain.ItemRepository,
	inventoryRepository domain.InventoryRepository, timeout time.Duration) domain.InventoryUseCase {
	return &inventoryUseCase{itemRepository: itemRepository, inventoryRepository: inventoryRepository, timeout: timeout}
}

func (i *inventoryUseCase) GetAll(ctx context.Context, count int, offset int,
	filter domain.InventorySpecification) ([]domain.InventoryItem, error) {
	c, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	inventoryItems, err := i.inventoryRepository.GetAll(c, count, offset, filter)
	if err != nil {
		return nil, err
	}
	if len(inventoryItems) == 0 {
		return nil, errors.NewNotFoundError("no items found matching the specification")
	}

	return i.fillItemDetails(c, inventoryItems)
}

func (i *inventoryUseCase) fillItemDetails(c context.Context, inventoryItems []domain.InventoryItem) ([]domain.InventoryItem, error) {
	group, ctx := errgroup.WithContext(c)

	inventoryMap := map[string]domain.Item{}

	for _, inventory := range inventoryItems {
		inventoryMap[inventory.Item.ID] = domain.Item{}
	}

	itemChan := make(chan domain.Item)
	for itemID := range inventoryMap {
		itemID := itemID
		group.Go(func() error {
			res, err := i.itemRepository.GetOne(ctx, itemID)
			if err != nil {
				return err
			}
			itemChan <- *res
			return nil
		})
	}

	go func() {
		err := group.Wait()
		if err != nil {
			log.Println(err.Error())
			return
		}
		close(itemChan)
	}()

	for item := range itemChan {
		if item != (domain.Item{}) {
			inventoryMap[item.ID] = item
		}
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}

	for index, inventoryItem := range inventoryItems {
		if item, ok := inventoryMap[inventoryItem.Item.ID]; ok {
			inventoryItems[index].Item = item
		}
	}
	return inventoryItems, nil
}

func (i *inventoryUseCase) GetInventoryForItem(ctx context.Context, itemID string) (*domain.InventoryItem, error) {
	c, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	inventoryItem, err := i.inventoryRepository.GetInventoryForItem(c, itemID)
	if err != nil {
		return nil, err
	}
	if inventoryItem == nil || inventoryItem.ID == "" {
		return nil, errors.NewNotFoundError("no item found matching the specification")
	}

	item, err := i.itemRepository.GetOne(c, inventoryItem.Item.ID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		log.Println("error getting the item for inventory")
	} else {
		inventoryItem.Item = *item
	}

	return inventoryItem, nil
}

func (i *inventoryUseCase) UpdateInventoryItem(ctx context.Context, inventory *domain.InventoryItem) (*domain.InventoryItem, error) {
	c, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	// no inventory already exists
	if inventory.ID == "" {
		var existingItem *domain.Item
		var err error
		// check if item has id or is new
		if inventory.Item.ID != "" {
			existingItem, err = i.itemRepository.GetOne(c, inventory.Item.ID)
			if err != nil {
				return nil, err
			}
			if existingItem == nil || existingItem.ID == "" {
				return nil, errors.NewBadRequestError("no item with such ID exists")
			}
		}
		// if no such item exists then create one and save it
		if existingItem == nil || existingItem.ID == "" {
			// create an item
			inventory.Item.ID = uuid.NewString()
			inventory.Item.CreatedAt = time.Now()
			_, err = i.itemRepository.Save(c, &inventory.Item)
			if err != nil {
				return nil, err
			}
		}
	}

	// this allows us to increase the number if item already exists in inventory instead of creating another
	// entry for the same item
	inv, err := i.inventoryRepository.GetInventoryForItem(c, inventory.Item.ID)
	if err != nil {
		return nil, err
	}
	if inventory.ID == "" {
		if (*inv).ID == "" {
			inventory.ID = uuid.NewString()
			return i.inventoryRepository.Save(c, inventory)
		} else {
			inventory.ID = inv.ID
		}
	}
	// ensure we aren't trying to change the item. why must this ever happen?
	if inv.ID != inventory.ID {
		log.Println("error", inv.ID, inventory.ID)
		return nil, errors.NewBadRequestError("invalid request. Can't change the item while updating")
	}
	if inventory.Quantity < 0 {
		return nil, errors.NewBadRequestError("invalid request. Quantity can't be less than 0")
	}
	inventory.UpdatedAt = time.Now()
	var updated *domain.InventoryItem
	updated, err = i.inventoryRepository.Edit(c, inventory)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (i *inventoryUseCase) DeleteItem(ctx context.Context, id string) error {
	c, cancel := context.WithTimeout(ctx, i.timeout)
	defer cancel()

	inventory, err := i.inventoryRepository.GetByID(c, id)
	if err != nil {
		log.Println(err)
		return errors.NewInternalServerError(err.Error())
	}
	if inventory == nil || inventory.ID == "" {
		return errors.NewNotFoundError("no such item found")
	}
	err = i.inventoryRepository.DeleteItem(c, inventory.ID)
	if err != nil {
		log.Println(err.Error())
		return errors.NewInternalServerError(err.Error())
	}

	// this is a design choice. Perhaps a bit iffy. In real life, it'd depend on what the client wants
	return i.itemRepository.Delete(c, inventory.Item.ID)
}

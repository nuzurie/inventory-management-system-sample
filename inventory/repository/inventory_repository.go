package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nuzurie/shopify/domain"
	"github.com/nuzurie/shopify/utils/errors"
	"log"
	"time"
)

type inventoryRepository struct {
	db *pgxpool.Pool
}

const (
	createInventoryTable = `CREATE TABLE IF NOT EXISTS inventory (
	id text PRIMARY KEY,
	quantity int,
	updated_at timestamp without time zone,
	item_id text,
	FOREIGN KEY (item_id)
	REFERENCES item(id)
	)`
	getInventoryForItemID = `SELECT id, quantity, updated_at, item_id WHERE item_id=$1`
	getAll                = `SELECT id, quantity, updated_at WHERE item_id IN %s AND %s LIMIT $1 OFFSET $2`
	getByID				  = `SELECT id, quantity, updated_at, item_id FROM public.inventory WHERE id=$1`
	save                  = `INSERT INTO public.inventory id, quantity, updated_at, item_id
			VALUES ($1, $2, $3, $4)`
	update        = `UPDATE public.inventory SET quantity=$2, updated_at=$3 WHERE id=$1`
	deleteForID = `DELETE FROM public.inventory WHERE id=$1`
)

func NewInventoryRepository(db *pgxpool.Pool) (domain.InventoryRepository, error) {
	log.Println("Creating item table")
	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
		//log.Fatalln("Unable to create the table. ", err.Error())
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), createInventoryTable)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &inventoryRepository{db: db}, nil
}

func (i *inventoryRepository) GetInventoryForItem(ctx context.Context, itemID string) (*domain.InventoryItem, error) {
	rows, err := i.db.Query(ctx, getInventoryForItemID, itemID)
	if err != nil {
		err = errors.NewInternalServerError(err.Error())
		return nil, err
	}
	defer rows.Close()

	var inventory domain.InventoryItem
	for rows.Next() {
		value, err := rows.Values()
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}

		inventory.ID = value[0].(string)
		inventory.Quantity = value[1].(int)
		inventory.UpdatedAt = value[2].(time.Time)
		inventory.Item.ID = value[3].(string)
	}

	return &inventory, nil
}

func (i *inventoryRepository) GetByID(ctx context.Context, id string) (*domain.InventoryItem, error) {
	rows, err := i.db.Query(ctx, getByID, id)
	if err != nil {
		err = errors.NewInternalServerError(err.Error())
		return nil, err
	}
	defer rows.Close()

	var inventory domain.InventoryItem
	for rows.Next() {
		value, err := rows.Values()
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}

		inventory.ID = value[0].(string)
		inventory.Quantity = value[1].(int)
		inventory.UpdatedAt = value[2].(time.Time)
		inventory.Item.ID = value[3].(string)
	}

	return &inventory, nil
}

func (i *inventoryRepository) GetAll(ctx context.Context, count int, offset int, filter domain.InventorySpecification) ([]domain.InventoryItem, error) {
	rows, err := i.db.Query(ctx, fmt.Sprintf(getAll, filter.ItemFilterQuery(), filter.FilterQuery()), count, offset)
	defer rows.Close()
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	var inventoryItems []domain.InventoryItem
	for rows.Next() {
		value, err := rows.Values()
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		var inventory domain.InventoryItem
		inventory.ID = value[0].(string)
		inventory.Quantity = value[1].(int)
		inventory.UpdatedAt = value[2].(time.Time)
		inventory.Item.ID = value[3].(string)

		inventoryItems = append(inventoryItems, inventory)
	}

	return inventoryItems, nil
}

func (i *inventoryRepository) Save(ctx context.Context, inventoryItem *domain.InventoryItem) (*domain.InventoryItem, error) {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, save, inventoryItem.ID, inventoryItem.Quantity, inventoryItem.UpdatedAt, inventoryItem.Item.ID)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	return inventoryItem, nil
}

func (i *inventoryRepository) Edit(ctx context.Context, inventoryItem *domain.InventoryItem) (*domain.InventoryItem, error) {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, update, inventoryItem.ID, inventoryItem.Quantity, inventoryItem.UpdatedAt, inventoryItem.Item.ID)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	return inventoryItem, nil
}

func (i *inventoryRepository) DeleteItem(ctx context.Context, id string) error {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, deleteForID, id)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

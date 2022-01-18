package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nuzurie/shopify/domain"
	"github.com/nuzurie/shopify/utils/errors"
	"log"
	"time"
)

type itemRepository struct {
	db *pgxpool.Pool
}

/*
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
 */

const (
	createItemTable = `CREATE TABLE IF NOT EXISTS item (
	id text PRIMARY KEY,
	name text NOT NULL,
	description text,
	price float,
	created_at timestamp without time zone,
	updated_at timestamp without time zone
	);`
	getByID = `SELECT id, name, description, price, created_at, updated_at FROM public.item WHERE id=$1`
	getAll = `SELECT id, name, description, price, created_at, updated_at FROM public.item WHERE $1 LIMIT $2 OFFSET $3`
	save = `INSERT INTO public.item(id, name, description, price, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6)`
	update = `UPDATE public.item
	SET name=$2, description=$3, price=$4, updated_at=$5
	WHERE id=$1;`
	deleteByID = `DELETE FROM public.item WHERE id=$1`
)

func NewItemRepository(db *pgxpool.Pool) (domain.ItemRepository, error) {
	log.Println("Creating item table")
	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
		//log.Fatalln("Unable to create the table. ", err.Error())
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), createItemTable)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &itemRepository{db: db}, nil
}

func (i itemRepository) GetAll(ctx context.Context, count int, offset int, filter domain.Specification) ([]domain.Item, error) {
	rows, err := i.db.Query(ctx, getAll, filter, count, offset)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	var items []domain.Item
	for rows.Next() {
		value, err := rows.Values()
		if err != nil {
			return nil, errors.NewInternalServerError(err.Error())
		}

		var item domain.Item
		item.ID = value[0].(string)
		item.Name = value[1].(string)
		item.Description = value[2].(string)
		item.Price = value[3].(float64)
		item.CreatedAt = value[4].(time.Time)
		item.UpdatedAt = value[5].(time.Time)

		items = append(items, item)
	}

	return items, nil
}

func (i itemRepository) GetOne(ctx context.Context, id string) (*domain.Item, error) {
	rows, err := i.db.Query(ctx, getByID, id)
	if err != nil {
		err = errors.NewInternalServerError(err.Error())
		return nil, err
	}
	defer rows.Close()

	var item domain.Item
	for rows.Next() {
		value, err := rows.Values()
		if err != nil {
			err = errors.NewInternalServerError(err.Error())
			return nil, err
		}

		item.ID = value[0].(string)
		item.Name = value[1].(string)
		item.Description = value[2].(string)
		item.Price = value[3].(float64)
		item.CreatedAt = value[4].(time.Time)
		item.UpdatedAt = value[5].(time.Time)
	}

	return &item, nil
}

func (i itemRepository) Save(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, save, item.ID, item.Name, item.Description, item.Price, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	return item, nil
}

func (i itemRepository) Edit(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, update, item.ID, item.Name, item.Description, item.Price, item.UpdatedAt)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	return item, nil
}

func (i itemRepository) Delete(ctx context.Context, id string) error {
	tx, err := i.db.Begin(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, deleteByID, id)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

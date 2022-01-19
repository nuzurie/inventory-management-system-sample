package specification

import (
	"fmt"
	"github.com/nuzurie/shopify/domain"
)

type InventorySpecification struct {
	ItemSpecification domain.Specification
	postgresQuery     string
}

func NewInventorySpecification(minQuantity, maxQuantity int,
	itemSpecification domain.Specification) domain.InventorySpecification {
	var minQuantityQuery, maxQuantityQuery string

	minQuantityQuery = fmt.Sprintf("quantity>=%d", minQuantity)
	if maxQuantity == -1 {
		maxQuantityQuery = "1=1"
	} else {
		maxQuantityQuery = fmt.Sprintf("quantity<=%d", maxQuantity)
	}

	query := fmt.Sprintf("%s AND %s", minQuantityQuery, maxQuantityQuery)

	return InventorySpecification{ItemSpecification: itemSpecification, postgresQuery: query}
}

func (i InventorySpecification) FilterQuery() string {
	return i.postgresQuery
}

func (i InventorySpecification) ItemFilterQuery() string {
	return i.ItemSpecification.FilterQuery()
}

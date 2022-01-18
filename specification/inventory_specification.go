package specification

import (
	"fmt"
	"github.com/nuzurie/shopify/domain"
)

type InventorySpecification struct {
	ItemSpecification domain.Specification
	postgresQuery     string
}

func NewInventorySpecification(minQuantity, maxQuantity int, itemSpecification domain.Specification) domain.Specification {
	var minQuantityQuery, maxQuntityQuery string

	minQuantityQuery = fmt.Sprintf("quantity>=%d", minQuantity)
	if maxQuantity == -1 {
		maxQuntityQuery = "1=1"
	} else {
		maxQuntityQuery = fmt.Sprintf("quantity<=%d", maxQuantity)
	}

	query := fmt.Sprintf("%s AND %s AND %s", itemSpecification.FilterQuery(), minQuantityQuery, maxQuntityQuery)

	return InventorySpecification{ItemSpecification: itemSpecification, postgresQuery: query}
}

func (i InventorySpecification) FilterQuery() string {
	return i.postgresQuery
}

package specification

import (
	"fmt"
	"github.com/nuzurie/shopify/domain"
)

type ItemSpecification struct {
	postgresQuery string
}

func NewItemSpecification(name, description string,
	minPrice, maxPrice float64) domain.Specification {
	var nameQuery, descriptionQuery, minPriceQuery, maxPriceQuery string
	if name == "" {
		nameQuery = "1=1"
	} else {
		nameQuery = fmt.Sprintf("name ILIKE %%%s%%", name)
	}

	if description == "" {
		descriptionQuery = "1=1"
	} else {
		descriptionQuery = fmt.Sprintf("description ILIKE '%%%s%%'", description)
	}

	minPriceQuery = fmt.Sprintf("price>=%f", minPrice)
	if maxPrice == -1 {
		maxPriceQuery = "1=1"
	} else {
		maxPriceQuery = fmt.Sprintf("price<=%f", maxPrice)
	}

	query := fmt.Sprintf("%s AND %s AND %s AND %s", nameQuery, descriptionQuery, minPriceQuery, maxPriceQuery)
	return ItemSpecification{postgresQuery: query}
}

func (i ItemSpecification) FilterQuery() string {
	return i.postgresQuery
}

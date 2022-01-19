package domain

type Specification interface {
	FilterQuery() string
}

type InventorySpecification interface {
	Specification
	ItemFilterQuery() string
}

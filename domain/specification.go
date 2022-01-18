package domain

type Specification interface {
	FilterQuery() string
}

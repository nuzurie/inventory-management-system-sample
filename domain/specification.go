package domain

type Specification interface {
	PostgresQuery() string
}

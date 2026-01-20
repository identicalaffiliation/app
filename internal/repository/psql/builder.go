package psql

import "github.com/Masterminds/squirrel"

type Builder struct {
	Builder squirrel.StatementBuilderType
}

func NewQueryBuilder() *Builder {
	return &Builder{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

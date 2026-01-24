package psql

import "github.com/Masterminds/squirrel"

type builder struct {
	Builder squirrel.StatementBuilderType
}

func NewQueryBuilder() *builder {
	return &builder{
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

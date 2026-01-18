package psql

import "github.com/Masterminds/squirrel"

type Builder struct {
	builder squirrel.StatementBuilderType
}

func NewQueryBuilder() *Builder {
	return &Builder{
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

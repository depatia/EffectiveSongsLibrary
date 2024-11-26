package tools

import (
	"fmt"
	"strings"
)

const (
	DEFAULT_OFFSET = 0
	DEFAULT_LIMIT  = 10
)

type ConditionQueryBuilder struct {
	builder strings.Builder
	args    []any

	wasWhere bool
}

func (b *ConditionQueryBuilder) Where(column string, arg any) *ConditionQueryBuilder {
	b.builder.WriteRune(' ')
	if !b.wasWhere {
		b.builder.WriteString("WHERE")
		b.wasWhere = true
	} else {
		b.builder.WriteString("AND")
	}
	b.builder.WriteString(fmt.Sprintf(" \"%s\" = $%d", column, len(b.args)+1))
	b.args = append(b.args, arg)
	return b
}

type SelectQueryBuilder struct {
	ConditionQueryBuilder
}

func (b *SelectQueryBuilder) Select(table string) *SelectQueryBuilder {
	b.builder.WriteString(fmt.Sprintf("SELECT * FROM %s", table))
	return b
}

func (b *SelectQueryBuilder) Offset(offset int) *SelectQueryBuilder {
	if offset == 0 {
		b.builder.WriteString(fmt.Sprintf(" OFFSET %d", DEFAULT_OFFSET))
	} else {
		b.builder.WriteString(fmt.Sprintf(" OFFSET %d", offset))
	}
	return b
}

func (b *SelectQueryBuilder) Limit(limit int) *SelectQueryBuilder {
	if limit == 0 {
		b.builder.WriteString(fmt.Sprintf(" LIMIT %d", DEFAULT_LIMIT))
	} else {
		b.builder.WriteString(fmt.Sprintf(" LIMIT %d", limit))
	}
	return b
}

func (b *SelectQueryBuilder) String() string {
	return b.builder.String()
}

func (b *SelectQueryBuilder) Args() []any {
	return b.args
}

type UpdateQueryBuilder struct {
	ConditionQueryBuilder

	wasSet bool
}

func (b *UpdateQueryBuilder) Update(tableName string) *UpdateQueryBuilder {
	b.builder.WriteString(fmt.Sprintf("UPDATE %s", tableName))
	return b
}

func (b *UpdateQueryBuilder) Set(column string, arg any) *UpdateQueryBuilder {
	if !b.wasSet {
		b.builder.WriteString(" SET")
		b.wasSet = true
	} else {
		b.builder.WriteRune(',')
	}
	b.builder.WriteString(fmt.Sprintf(" \"%s\" = $%d", column, len(b.args)+1))
	b.args = append(b.args, arg)
	return b
}

func (b *UpdateQueryBuilder) Args() []any {
	return b.args
}

func (b *UpdateQueryBuilder) String() string {
	return b.builder.String()
}

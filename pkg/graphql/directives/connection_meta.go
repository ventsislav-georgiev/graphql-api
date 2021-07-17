package directives

import (
	"github.com/graph-gophers/graphql-go/types"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers"
)

type ConnectionMeta struct {
	Embedded     bool
	PrimaryKey   *string
	PrimaryValue *string
	ForeignKey   *string
}

func GetConnectionMetaFromDefinition(fieldDefinition types.FieldDefinition, defaultIDFieldName string) ConnectionMeta {
	connectionMeta := ConnectionMeta{
		Embedded: fieldDefinition.Parent != nil,
	}

	directives := fieldDefinition.Directives
	if directives == nil {
		return connectionMeta
	}

	meta := directives.Get(ConnectionDirectiveName)
	if meta == nil {
		return connectionMeta
	}

	connectionMeta.Embedded = false

	primaryKey, ok := meta.Arguments.Get("primaryKey")
	if ok && primaryKey != nil {
		connectionMeta.PrimaryKey = helpers.String(primaryKey.String())
	} else {
		connectionMeta.PrimaryKey = helpers.String(defaultIDFieldName)
	}

	foreignKey, ok := meta.Arguments.Get("foreignKey")
	if ok && foreignKey != nil {
		connectionMeta.ForeignKey = helpers.String(foreignKey.String())
	} else {
		connectionMeta.ForeignKey = helpers.String(defaultIDFieldName)
	}

	return connectionMeta
}

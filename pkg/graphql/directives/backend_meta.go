package directives

import (
	"strings"

	"github.com/graph-gophers/graphql-go/types"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers"
)

type BackendMeta struct {
	Product        *string
	CollectionName *string
	Method         *string
	Endpoint       *string
	IsList         bool
}

func GetBackendMetaFromDefinition(fieldDefinition types.FieldDefinition) BackendMeta {
	// Check if field return type is List
	list, isList := fieldDefinition.Type.(*types.List)

	// If field has @backend() directive
	if HasBackendDirective(fieldDefinition.Directives) {
		return GetBackendMeta(fieldDefinition.Directives, isList)
	}

	// If field return type is object
	if objectTypeDefinition, ok := fieldDefinition.Type.(*types.ObjectTypeDefinition); ok && HasBackendDirective(objectTypeDefinition.Directives) {
		return GetBackendMeta(objectTypeDefinition.Directives, isList)
	}

	if isList {
		// If field return type is list of objects
		if objectTypeDefinition, ok := list.OfType.(*types.ObjectTypeDefinition); ok && HasBackendDirective(objectTypeDefinition.Directives) {
			return GetBackendMeta(objectTypeDefinition.Directives, isList)
		}
	}

	if fieldDefinition.Parent == nil {
		return BackendMeta{}
	}

	// Follow parent field definition
	backendMeta := GetBackendMetaFromDefinition(*fieldDefinition.Parent)
	backendMeta.IsList = isList
	return backendMeta
}

func GetBackendMeta(directives types.DirectiveList, isList bool) BackendMeta {
	backendMeta := BackendMeta{IsList: isList}

	if directives == nil {
		return backendMeta
	}

	meta := directives.Get(BackendDirectiveName)
	if meta == nil {
		return backendMeta
	}

	product, ok := meta.Arguments.Get("product")
	if ok && product != nil {
		backendMeta.Product = helpers.String(product.String())
	}

	collectionName, ok := meta.Arguments.Get("collection")
	if ok && collectionName != nil {
		backendMeta.CollectionName = helpers.String(collectionName.String())
	}

	method, ok := meta.Arguments.Get("method")
	if ok && method != nil {
		backendMeta.Method = helpers.String(method.String())
	}

	endpoint, ok := meta.Arguments.Get("endpoint")
	if ok && endpoint != nil {
		backendMeta.Endpoint = helpers.String(strings.ReplaceAll(endpoint.String(), "\"", ""))
	}

	return backendMeta
}

func HasBackendDirective(fieldDirectives types.DirectiveList) bool {
	if fieldDirectives == nil {
		return false
	}

	meta := fieldDirectives.Get(BackendDirectiveName)
	return meta != nil
}

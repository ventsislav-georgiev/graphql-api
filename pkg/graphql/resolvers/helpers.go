package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go/types"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/httpmethod"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/responsetype"
	"github.com/ventsislav-georgiev/graphql-api/pkg/graphql/directives"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers"
)

func scalar(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, args map[string]interface{}, backendResolver BackendResolver) types.DynamicResolver {
	metaDirectives := &directives.MetaDirectives{
		Backend:    directives.GetBackendMetaFromDefinition(fieldDefinition),
		Connection: directives.GetConnectionMetaFromDefinition(fieldDefinition, backends.DefaultKinveyIdentifierFieldName),
	}

	resolvedData := getResolvedData(pathSegment)
	var data interface{}
	var found bool

	if resolvedData != nil {
		data, found = resolvedData[fieldDefinition.Name]
	}

	if found {
		return getDynamicResult(data)
	} else if metaDirectives.Connection.Embedded {
		return &dynamicResolver{err: customErrorNotFound}
	}

	if resolvedData != nil && metaDirectives.Connection.PrimaryKey != nil {
		metaDirectives.Connection.PrimaryValue = helpers.GetString(resolvedData, *metaDirectives.Connection.PrimaryKey)
	}

	result, err := backendResolver.GetDataFromBackend(ctx, fieldDefinition, pathSegment, selectedFields, metaDirectives, args)
	if err != nil {
		return &dynamicResolver{err: &customError{err: err}}
	}

	return getDynamicResult(result)
}

func list(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, args map[string]interface{}, backendResolver BackendResolver) (*[]interface{}, error) {
	metaDirectives := &directives.MetaDirectives{
		Backend:    directives.GetBackendMetaFromDefinition(fieldDefinition),
		Connection: directives.GetConnectionMetaFromDefinition(fieldDefinition, backends.DefaultKinveyIdentifierFieldName),
	}

	resolvedData := getResolvedData(pathSegment)
	if metaDirectives.Connection.Embedded {
		if resolvedData == nil {
			return nil, customErrorNotFound
		}

		value, found := resolvedData[fieldDefinition.Name]
		if !found {
			return nil, customErrorNotFound
		}

		return getDynamicResultList(value)
	}

	result, err := backendResolver.GetDataFromBackend(ctx, fieldDefinition, pathSegment, selectedFields, metaDirectives, args)
	if err != nil {
		return nil, &customError{err: err}
	}

	return getDynamicResultList(result)
}

func getMethod(metaDirectives *directives.MetaDirectives) httpmethod.HttpMethod {
	method := httpmethod.GET
	if metaDirectives.Backend.Method != nil {
		method = httpmethod.FromString(*metaDirectives.Backend.Method)
	}
	return method
}

func getResponseType(fieldDefinition types.FieldDefinition) responsetype.ResponseType {
	responseType := responsetype.JSONObject
	if _, ok := fieldDefinition.Type.(*types.ScalarTypeDefinition); ok {
		responseType = responsetype.ScalarValue
	} else if _, ok := fieldDefinition.Type.(*types.ObjectTypeDefinition); ok {
		responseType = responsetype.JSONObject
	} else if _, ok := fieldDefinition.Type.(*types.List); ok {
		responseType = responsetype.JSONArrayOfObjects
	}
	return responseType
}

func getResolvedData(pathSegment types.PathSegment) map[string]interface{} {
	segment := &pathSegment
	for segment.Resolver.IsNil() && segment.Parent != nil {
		segment = segment.Parent
	}

	if segment.Resolver.IsNil() {
		return nil
	}

	parentData, ok := segment.Resolver.Interface().(*dynamicResolver)
	if !ok {
		return nil
	}

	return parentData.data
}

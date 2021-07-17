package resolvers

import (
	"context"
	"errors"
	"sync"

	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go/types"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/httpmethod"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/responsetype"
	"github.com/ventsislav-georgiev/graphql-api/pkg/graphql/directives"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers/constants"
)

type SitefinityResolver struct {
	Backend     backends.SitefinityBackendProvider
	DataLoaders *sync.Map
}

func (p *SitefinityResolver) Scalar(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, args map[string]interface{}) types.DynamicResolver {
	return scalar(ctx, fieldDefinition, pathSegment, selectedFields, args, p)
}

func (p *SitefinityResolver) List(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, args map[string]interface{}) (*[]interface{}, error) {
	return list(ctx, fieldDefinition, pathSegment, selectedFields, args, p)
}

func (p *SitefinityResolver) GetDataFromBackend(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, metaDirectives *directives.MetaDirectives, args map[string]interface{}) (interface{}, error) {
	var ok bool
	var endpoint string

	method := getMethod(metaDirectives)
	responseType := getResponseType(fieldDefinition)

	if metaDirectives.Backend.Endpoint != nil {
		endpoint = *metaDirectives.Backend.Endpoint
	} else if metaDirectives.Backend.CollectionName != nil {
		collectionName := *metaDirectives.Backend.CollectionName
		endpoint = "/api/default/" + collectionName
	} else {
		return nil, errors.New("failed to resolve backend endpoint")
	}

	id := helpers.GetString(args, constants.IdParamName)
	if id == nil && metaDirectives.Connection.PrimaryValue != nil {
		id = metaDirectives.Connection.PrimaryValue
		responseType = responsetype.JSONObject
	}

	var data map[string]interface{}
	if dataArg, found := args[constants.DataParamName]; found {
		if data, ok = dataArg.(map[string]interface{}); !ok {
			return nil, errors.New("failed to cast input data as object")
		}
	}

	filter := helpers.GetString(args, constants.ODataFilterParamName)
	sort := helpers.GetString(args, constants.ODataSortParamName)
	expand := helpers.AsString(nil)

	var objectTypeDefinition *types.ObjectTypeDefinition
	if list, isList := fieldDefinition.Type.(*types.List); isList {
		objectTypeDefinition, _ = list.OfType.(*types.ObjectTypeDefinition)
	} else {
		objectTypeDefinition, _ = fieldDefinition.Type.(*types.ObjectTypeDefinition)
	}

	if objectTypeDefinition != nil && method == httpmethod.GET {
		fields := ""
		for _, f := range objectTypeDefinition.Fields {
			if f.Directives.Get(directives.ConnectionDirectiveName) != nil && helpers.SelectedFieldsContains(selectedFields.Fields, f.Name) {
				fields += f.Name + ","
			}
		}

		if !helpers.IsEmpty(fields) {
			expand = &fields
		}
	}

	if method != httpmethod.GET {
		if id != nil {
			endpoint += "/" + *id
		}
		return p.Backend.Request(endpoint, method, responseType, data, filter, sort, expand)
	}

	ctx = backends.CtxSetSitefinityBackendProvider(ctx, p.Backend)
	ctx = backends.CtxSetHttpEndpoint(ctx, endpoint)
	ctx = backends.CtxSetHttpMethod(ctx, method)
	ctx = backends.CtxSetResponseType(ctx, responseType)
	ctx = backends.CtxSetHttpData(ctx, data)
	ctx = backends.CtxSetFilter(ctx, filter)
	ctx = backends.CtxSetSort(ctx, sort)
	ctx = backends.CtxSetExpand(ctx, expand)

	var thunk dataloader.Thunk
	if id != nil {
		loader := helpers.GetOrAddLoader(p.DataLoaders, backends.SitefinityLoadBatchFn, &endpoint, helpers.String("withId"))
		thunk = loader.Load(ctx, dataloader.StringKey(*id))
	} else {
		loader := helpers.GetOrAddLoader(p.DataLoaders, backends.SitefinityLoadManyBatchFn, &endpoint, filter, sort, expand)
		thunk = loader.Load(ctx, dataloader.StringKey(endpoint))
	}

	result, err := thunk()
	return result, err
}

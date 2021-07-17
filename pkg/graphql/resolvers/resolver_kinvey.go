package resolvers

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go/types"
	"github.com/intel/rsp-sw-toolkit-im-suite-go-odata/odata"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/httpmethod"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/responsetype"
	"github.com/ventsislav-georgiev/graphql-api/pkg/graphql/directives"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers/constants"
)

type KinveyResolver struct {
	Backend     backends.KinveyBackendProvider
	DataLoaders *sync.Map
}

func (p *KinveyResolver) Scalar(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, args map[string]interface{}) types.DynamicResolver {
	return scalar(ctx, fieldDefinition, pathSegment, selectedFields, args, p)
}

func (p *KinveyResolver) List(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, args map[string]interface{}) (*[]interface{}, error) {
	return list(ctx, fieldDefinition, pathSegment, selectedFields, args, p)
}

func (p *KinveyResolver) GetDataFromBackend(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, metaDirectives *directives.MetaDirectives, args map[string]interface{}) (interface{}, error) {
	var ok bool
	var err error
	var endpoint string

	method := getMethod(metaDirectives)
	responseType := getResponseType(fieldDefinition)

	if metaDirectives.Backend.Endpoint != nil {
		endpoint = *metaDirectives.Backend.Endpoint
		endpoint = strings.ReplaceAll(endpoint, ":kid", p.Backend.KinveyID)
	} else if metaDirectives.Backend.CollectionName != nil {
		collectionName := *metaDirectives.Backend.CollectionName
		endpoint = "/appdata/" + p.Backend.KinveyID + "/" + collectionName
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

	var filter *string
	if odataString := helpers.GetString(args, constants.ODataFilterParamName); odataString != nil {
		filterQuery, err := odata.ParseODataFilter(*odataString)
		if err != nil {
			return nil, err
		}
		filterMap, err := odata.ParseODataFilterForMongo(filterQuery)
		if err != nil {
			return nil, err
		}
		jsonBytes, err := json.Marshal(filterMap)
		if err != nil {
			return nil, err
		}
		filter = helpers.String(string(jsonBytes))
	}

	var sort *string
	if odataString := helpers.GetString(args, constants.ODataSortParamName); odataString != nil {
		sortMap, err := odata.ParseODataOrderBy(*odataString)
		if err != nil {
			return nil, err
		}
		jsonBytes, err := json.Marshal(sortMap)
		if err != nil {
			return nil, err
		}
		sort = helpers.String(string(jsonBytes))
	}

	if method != httpmethod.GET {
		if id != nil {
			endpoint += "/" + *id
		}
		return p.Backend.Request(endpoint, method, responseType, data, filter, sort)
	}

	ctx = backends.CtxSetKinveyBackendProvider(ctx, p.Backend)
	ctx = backends.CtxSetHttpEndpoint(ctx, endpoint)
	ctx = backends.CtxSetHttpMethod(ctx, method)
	ctx = backends.CtxSetResponseType(ctx, responseType)
	ctx = backends.CtxSetHttpData(ctx, data)
	ctx = backends.CtxSetFilter(ctx, filter)
	ctx = backends.CtxSetSort(ctx, sort)

	var thunk dataloader.Thunk
	if id != nil {
		loader := helpers.GetOrAddLoader(p.DataLoaders, backends.KinveyLoadBatchFn, &endpoint, helpers.String("withId"))
		thunk = loader.Load(ctx, dataloader.StringKey(*id))
	} else {
		loader := helpers.GetOrAddLoader(p.DataLoaders, backends.KinveyLoadManyBatchFn, &endpoint, filter, sort)
		thunk = loader.Load(ctx, dataloader.StringKey(endpoint))
	}

	result, err := thunk()
	return result, err
}

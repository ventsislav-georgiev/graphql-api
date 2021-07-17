package backends

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/graph-gophers/dataloader"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/responsetype"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers"
)

const (
	DefaultSitefinityIdentifierFieldName = "Id"
)

func SitefinityLoadBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	var result interface{}
	var err error

	results := make([]*dataloader.Result, len(keys))
	fail := func() {
		log.Fatal("context values not passed properly")
	}

	backend, ok := CtxGetSitefinityBackendProvider(ctx)
	if !ok {
		fail()
	}
	endpoint, ok := CtxGetHttpEndpoint(ctx)
	if !ok {
		fail()
	}
	method, ok := CtxGetHttpMethod(ctx)
	if !ok {
		fail()
	}

	batchFilter := DefaultSitefinityIdentifierFieldName + " in ('" + strings.Join(keys.Keys(), "', '") + "')"
	result, err = backend.Request(endpoint, method, responsetype.JSONArrayOfObjects, nil, &batchFilter, nil, nil)

	var arrayOfObjects []map[string]interface{}
	if err == nil {
		if arrayOfObjects, ok = result.([]map[string]interface{}); !ok {
			err = errors.New("invalid batch response for " + endpoint)
		}
	}

	if err != nil {
		for i := range keys {
			results[i] = &dataloader.Result{Error: err}
		}
		return results
	}

	for i, key := range keys.Keys() {
		for _, obj := range arrayOfObjects {
			if helpers.GetStringOrEmpty(obj, DefaultSitefinityIdentifierFieldName) == key {
				results[i] = &dataloader.Result{Data: obj}
				break
			}
		}
	}

	return results
}

func SitefinityLoadManyBatchFn(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	fail := func() {
		log.Fatal("context values not passed properly")
	}

	results := make([]*dataloader.Result, len(keys))
	backend, ok := CtxGetSitefinityBackendProvider(ctx)
	if !ok {
		fail()
	}
	endpoint, ok := CtxGetHttpEndpoint(ctx)
	if !ok {
		fail()
	}
	responsetype, ok := CtxGetResponseType(ctx)
	if !ok {
		fail()
	}
	method, ok := CtxGetHttpMethod(ctx)
	if !ok {
		fail()
	}
	filter, ok := CtxGetFilter(ctx)
	if !ok {
		fail()
	}
	sort, ok := CtxGetSort(ctx)
	if !ok {
		fail()
	}
	expand, ok := CtxGetExpand(ctx)
	if !ok {
		fail()
	}

	result, err := backend.Request(endpoint, method, responsetype, nil, filter, sort, expand)

	for i := range keys.Keys() {
		if err != nil {
			results[i] = &dataloader.Result{Error: err}
		} else {
			results[i] = &dataloader.Result{Data: result}
		}
	}

	return results
}

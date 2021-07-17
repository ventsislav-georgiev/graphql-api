package backends

import (
	"context"

	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/httpmethod"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends/responsetype"
)

type contextKey string

var (
	ctxKinveyBackendProviderKey     = contextKey("kinveyBackendProvider")
	ctxSitefinityBackendProviderKey = contextKey("sitefinityBackendProvider")
	ctxHttpEndpointKey              = contextKey("httpEndpoint")
	ctxHttpMethodKey                = contextKey("httpMethod")
	ctxHttpDataKey                  = contextKey("httpData")
	ctxResponseTypeKey              = contextKey("responseType")
	ctxFilterKey                    = contextKey("filter")
	ctxSortKey                      = contextKey("sort")
	ctxExpandKey                    = contextKey("expand")
)

func CtxSetKinveyBackendProvider(ctx context.Context, value KinveyBackendProvider) context.Context {
	return context.WithValue(ctx, ctxKinveyBackendProviderKey, value)
}

func CtxGetKinveyBackendProvider(ctx context.Context) (KinveyBackendProvider, bool) {
	val, ok := ctx.Value(ctxKinveyBackendProviderKey).(KinveyBackendProvider)
	return val, ok
}

func CtxSetSitefinityBackendProvider(ctx context.Context, value SitefinityBackendProvider) context.Context {
	return context.WithValue(ctx, ctxSitefinityBackendProviderKey, value)
}

func CtxGetSitefinityBackendProvider(ctx context.Context) (SitefinityBackendProvider, bool) {
	val, ok := ctx.Value(ctxSitefinityBackendProviderKey).(SitefinityBackendProvider)
	return val, ok
}

func CtxSetHttpEndpoint(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, ctxHttpEndpointKey, value)
}

func CtxGetHttpEndpoint(ctx context.Context) (string, bool) {
	val, ok := ctx.Value(ctxHttpEndpointKey).(string)
	return val, ok
}

func CtxSetHttpMethod(ctx context.Context, value httpmethod.HttpMethod) context.Context {
	return context.WithValue(ctx, ctxHttpMethodKey, value)
}

func CtxGetHttpMethod(ctx context.Context) (httpmethod.HttpMethod, bool) {
	val, ok := ctx.Value(ctxHttpMethodKey).(httpmethod.HttpMethod)
	return val, ok
}

func CtxSetHttpData(ctx context.Context, value map[string]interface{}) context.Context {
	return context.WithValue(ctx, ctxHttpDataKey, value)
}

func CtxGetHttpData(ctx context.Context) (map[string]interface{}, bool) {
	val, ok := ctx.Value(ctxHttpDataKey).(map[string]interface{})
	return val, ok
}

func CtxSetResponseType(ctx context.Context, value responsetype.ResponseType) context.Context {
	return context.WithValue(ctx, ctxResponseTypeKey, value)
}

func CtxGetResponseType(ctx context.Context) (responsetype.ResponseType, bool) {
	val, ok := ctx.Value(ctxResponseTypeKey).(responsetype.ResponseType)
	return val, ok
}

func CtxSetFilter(ctx context.Context, value *string) context.Context {
	return context.WithValue(ctx, ctxFilterKey, value)
}

func CtxGetFilter(ctx context.Context) (*string, bool) {
	val, ok := ctx.Value(ctxFilterKey).(*string)
	return val, ok
}

func CtxSetSort(ctx context.Context, value *string) context.Context {
	return context.WithValue(ctx, ctxSortKey, value)
}

func CtxGetSort(ctx context.Context) (*string, bool) {
	val, ok := ctx.Value(ctxSortKey).(*string)
	return val, ok
}

func CtxSetExpand(ctx context.Context, value *string) context.Context {
	return context.WithValue(ctx, ctxExpandKey, value)
}

func CtxGetExpand(ctx context.Context) (*string, bool) {
	val, ok := ctx.Value(ctxExpandKey).(*string)
	return val, ok
}

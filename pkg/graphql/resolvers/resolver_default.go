package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go/types"
	"github.com/ventsislav-georgiev/graphql-api/pkg/graphql/directives"
)

const (
	scalarResultMethodName = "Scalar"
	listResultMethodName   = "List"
)

type BackendResolver interface {
	Scalar(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, args map[string]interface{}) types.DynamicResolver
	List(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, args map[string]interface{}) (*[]interface{}, error)
	GetDataFromBackend(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, selectedFields types.SelectedFields, metaDirectives *directives.MetaDirectives, args map[string]interface{}) (interface{}, error)
}

type DefaultResolversProvider struct {
	KinveyResolver     BackendResolver
	SitefinityResolver BackendResolver
}

func (p *DefaultResolversProvider) GetResolver(fieldDefinition types.FieldDefinition) *types.ResolverInfo {
	backendMeta := directives.GetBackendMetaFromDefinition(fieldDefinition)
	if backendMeta.Product == nil {
		return nil
	}

	methodName := scalarResultMethodName
	if backendMeta.IsList {
		methodName = listResultMethodName
	}

	name := *backendMeta.Product
	switch name {
	case "kinvey":
		return &types.ResolverInfo{
			Resolver:   p.KinveyResolver,
			MethodName: methodName,
		}
	case "sitefinity":
		return &types.ResolverInfo{
			Resolver:   p.SitefinityResolver,
			MethodName: methodName,
		}
	default:
		return nil
	}
}

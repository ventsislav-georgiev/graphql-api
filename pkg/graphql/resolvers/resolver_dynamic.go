package resolvers

import (
	"context"
	"reflect"

	"github.com/graph-gophers/graphql-go/types"
)

type dynamicResolver struct {
	data           map[string]interface{}
	err            resolverError
	hasScalarValue bool
	scalarValue    interface{}
}

func (r *dynamicResolver) Resolve(ctx context.Context, fieldDefinition types.FieldDefinition, pathSegment types.PathSegment, args map[string]interface{}) (interface{}, error) {
	if r.err != nil {
		return nil, r.err
	}

	if r.hasScalarValue {
		return r.scalarValue, nil
	}

	if r.data == nil {
		return nil, customErrorNotFound
	}

	value, found := r.data[fieldDefinition.Name]
	if !found {
		return nil, customErrorNotFound
	}

	return value, nil
}

func (r *dynamicResolver) HasScalarValue() bool {
	if r.err != nil {
		return true
	}

	return r.hasScalarValue
}

func getDynamicResult(result interface{}) types.DynamicResolver {
	if result == nil {
		return &dynamicResolver{err: customErrorNotFound}
	}

	dictResult, ok := result.(map[string]interface{})
	if ok {
		return &dynamicResolver{data: dictResult}
	}

	return &dynamicResolver{hasScalarValue: true, scalarValue: result}
}

func getDynamicResultList(result interface{}) (*[]interface{}, error) {
	if result == nil {
		return nil, customErrorNotFound
	}

	resultType := reflect.TypeOf(result)
	resultValue := reflect.ValueOf(result)
	if resultType.Kind() == reflect.Ptr {
		resultType = resultType.Elem()
		resultValue = resultValue.Elem()
	}

	var dynamicResults []interface{}
	switch resultType.Kind() {
	case reflect.Slice:
		slice := resultValue
		dynamicResults = make([]interface{}, resultValue.Len())
		for i := 0; i < slice.Len(); i++ {
			value := slice.Index(i).Interface()
			dictValue, ok := value.(map[string]interface{})
			if ok {
				dynamicResults[i] = &dynamicResolver{data: dictValue}
			} else {
				dynamicResults[i] = value
			}
		}
	}

	return &dynamicResults, nil
}

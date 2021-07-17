package helpers

import (
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go/types"
)

type odataArrayOfObjectsResp struct {
	Value []map[string]interface{} `json:"value"`
}

func String(str string) *string {
	return &str
}

func AsString(str interface{}) *string {
	value, ok := str.(string)
	if !ok {
		return nil
	}

	return &value
}

func AsStringOrEmpty(str interface{}) string {
	value, ok := str.(string)
	if !ok {
		return ""
	}

	return value
}

func GetString(args map[string]interface{}, key string) *string {
	if value, ok := args[key]; ok {
		return AsString(value)
	}

	return nil
}

func GetStringOrEmpty(args map[string]interface{}, key string) string {
	if value, ok := args[key]; ok {
		return AsStringOrEmpty(value)
	}

	return ""
}

func LowercaseDictKeys(data map[string]interface{}) map[string]interface{} {
	for k, v := range data {
		lowercaseKey := strings.ToLower(k)
		if k != lowercaseKey {
			data[lowercaseKey] = v
			delete(data, k)
		}
	}
	return data
}

func JSONEncodeObject(data map[string]interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func JSONDecodeObject(data io.Reader) (map[string]interface{}, error) {
	jsonResult := make(map[string]interface{})
	err := json.NewDecoder(data).Decode(&jsonResult)
	if err != nil {
		return nil, err
	}
	return jsonResult, nil
}

func JSONDecodeODataArrayOfObjects(data io.Reader) ([]map[string]interface{}, error) {
	jsonResult := &odataArrayOfObjectsResp{}
	err := json.NewDecoder(data).Decode(&jsonResult)
	if err != nil {
		return nil, err
	}
	return jsonResult.Value, nil
}

func JSONDecodeArrayOfObjects(data io.Reader) ([]map[string]interface{}, error) {
	jsonResult := make([]map[string]interface{}, 0)
	err := json.NewDecoder(data).Decode(&jsonResult)
	if err != nil {
		return nil, err
	}
	return jsonResult, nil
}

func IsNullOrEmpty(str *string) bool {
	return str == nil || IsEmpty(*str)
}

func IsEmpty(str string) bool {
	return str == ""
}

func SelectedFieldsContains(selectedFields []*types.SelectedField, fieldName string) bool {
	for _, sf := range selectedFields {
		if sf.Name == fieldName {
			return true
		}
	}

	return false
}

func CreateKey(nullableKeyParts ...*string) string {
	keyParts := make([]string, 0)
	for _, part := range nullableKeyParts {
		if part != nil {
			keyParts = append(keyParts, *part)
		}
	}
	return strings.Join(keyParts, ".")
}

func GetOrAddLoader(dataLoaders *sync.Map, batchFunc dataloader.BatchFunc, loaderKeys ...*string) *dataloader.Loader {
	key := CreateKey(loaderKeys...)
	value, found := dataLoaders.Load(key)
	if found {
		return value.(*dataloader.Loader)
	}

	loader := dataloader.NewBatchedLoader(batchFunc, dataloader.WithWait(time.Millisecond*100))
	dataLoaders.Store(key, loader)
	return loader
}

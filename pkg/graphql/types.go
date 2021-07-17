package graphql

import (
	"encoding/json"
	"fmt"
)

type emptyStruct struct{}

type queryParams struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type batchParams []queryParams

func (r *batchParams) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("no bytes to unmarshal")
	}
	switch b[0] {
	case '{':
		return r.unmarshalSingle(b)
	case '[':
		return r.unmarshalMany(b)
	}

	err := r.unmarshalMany(b)
	if err != nil {
		return r.unmarshalSingle(b)
	}
	return nil
}

func (r *batchParams) unmarshalSingle(b []byte) error {
	var params queryParams
	err := json.Unmarshal(b, &params)
	*r = append(*r, params)
	return err
}

func (r *batchParams) unmarshalMany(b []byte) error {
	var batch []queryParams
	err := json.Unmarshal(b, &batch)
	*r = batch
	return err
}

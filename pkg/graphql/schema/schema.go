package schema

import "net/http"

type SchemaProvider interface {
	GetSchema(r *http.Request) *string
}

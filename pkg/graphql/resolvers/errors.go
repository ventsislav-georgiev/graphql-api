package resolvers

import "errors"

type resolverError interface {
	error
	Extensions() map[string]interface{}
}

type customError struct {
	err        error
	extensions map[string]interface{}
}

func (c customError) Error() string {
	return c.err.Error()
}

func (c customError) Extensions() map[string]interface{} {
	if c.extensions != nil {
		return c.extensions
	}

	return map[string]interface{}{}
}

var (
	errNotFound         = errors.New("data not found")
	customErrorNotFound = customError{err: errNotFound}
)

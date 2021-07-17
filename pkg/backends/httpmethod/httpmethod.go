package httpmethod

type HttpMethod int

const (
	GET HttpMethod = iota
	OPTIONS
	HEAD
	POST
	PUT
	PATCH
	DELETE
)

func (m HttpMethod) String() string {
	switch m {
	case GET:
		return "GET"
	case OPTIONS:
		return "OPTIONS"
	case HEAD:
		return "HEAD"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case PATCH:
		return "PATCH"
	case DELETE:
		return "DELETE"
	default:
		return ""
	}
}

func FromString(method string) HttpMethod {
	switch method {
	case "GET":
		return GET
	case "OPTIONS":
		return OPTIONS
	case "HEAD":
		return HEAD
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "PATCH":
		return PATCH
	case "DELETE":
		return DELETE
	default:
		return -1
	}
}

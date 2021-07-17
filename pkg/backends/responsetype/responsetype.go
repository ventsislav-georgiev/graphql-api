package responsetype

type ResponseType int

const (
	JSONObject ResponseType = iota
	JSONArrayOfObjects
	ScalarValue
)

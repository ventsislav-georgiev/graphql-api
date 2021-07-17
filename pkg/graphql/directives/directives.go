package directives

const (
	BackendDirectiveName    = "backend"
	ConnectionDirectiveName = "connection"
)

type MetaDirectives struct {
	Backend    BackendMeta
	Connection ConnectionMeta
}

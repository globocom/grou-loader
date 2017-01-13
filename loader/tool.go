package loader

type Tool interface {
	GetName() string
	Start(map[string]interface{})
	Check() bool
	Params() map[string]interface{}
}

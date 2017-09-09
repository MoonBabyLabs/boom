package base

type ServiceContract interface {
	Get(resource string) ModelContract
	DeleteResource(resource string) bool
}

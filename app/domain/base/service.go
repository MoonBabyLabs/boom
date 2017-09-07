package base

type ServiceContract interface {
	All() [3]int
	Get(resource string) ModelContract
	DeleteResource(resource string) bool
}

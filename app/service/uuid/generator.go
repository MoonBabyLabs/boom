package uuid

type Generator interface {
	New(hash string) string
	Init() Generator
}

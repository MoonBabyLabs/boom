package uuid

import (
	"github.com/satori/go.uuid"
)

type Default struct {

}

func (d Default) New(str string) string {
	return uuid.NewV5(uuid.NamespaceDNS, str).String()
}

func (d Default) Init() Generator {
	return d
}
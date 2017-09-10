package factory

import (
	"github.com/MoonBabyLabs/boom/app/domain/base"
)

type ServiceContract interface {
	run(domain base.ServiceContract)
}

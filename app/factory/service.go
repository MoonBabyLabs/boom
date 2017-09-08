package factory

import (
	"gigdub/app/domain/base"
)

type ServiceContract interface {
	run(domain base.ServiceContract)
}

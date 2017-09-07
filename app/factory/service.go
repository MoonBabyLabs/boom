package factory

import (
	"gigdubservice/app/domain/base"
)

type ServiceContract interface {
	run(domain base.ServiceContract)
}

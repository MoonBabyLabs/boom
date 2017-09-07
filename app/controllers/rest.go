package controllers

import (
	"github.com/revel/revel"
	"gigdubservice/app/domain/job"
	"gigdubservice/app/factory"
	"gigdubservice/app/domain/base"
	"errors"
	"log"
)


type Rest struct {
	*revel.Controller
}

type DataResponse struct {
	id int
}

type FactoryGenerator struct {
	service []factory.ServiceContract
}

type JobService struct {
	*job.Service
}


func (c Rest) Get(domain string, resource string) revel.Result {
	data := make(map[string]interface{})
	service, err := Generate(domain)
	newData := service.Get(resource)
	data["data"] = newData
	log.Print(newData)

	if err==nil  {
		return c.RenderJSON(newData)
	} else {
		data["error"] = err
		data["sucess"] = "false"
		log.Panic(err)
	}

	return c.RenderJSON(data);
}

func Generate(domainService string) (base.ServiceContract, error) {
	switch domainService {
	case "job":
		return new(job.Service), nil
	default:
		return nil, errors.New("Could not find domain item")
	}

}
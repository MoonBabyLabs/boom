package controllers

import (
	"github.com/revel/revel"
	"gigdub/app/domain/job"
	"gigdub/app/factory"
	"gigdub/app/domain/base"
	"errors"
	"log"
	"gigdub/app/provider"
	"encoding/json"
	"io/ioutil"
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
	}

	data["error"] = err
	data["sucess"] = "false"
	log.Panic(err)

	return c.RenderJSON(data);
}

func (c Rest) Post(domain string) revel.Result {
	model := base.Model{}
	item := make(map[string]interface{})
	content, _ := ioutil.ReadAll(c.Request.Body)
	json.Unmarshal([]byte(content), item)
	model.Domain = domain
	model.Datastore = provider.Db{}.Construct()
	model.Datastore.SetDomain(domain)
	model.Domain = domain
	model.Add(item)

	return c.RenderJSON(item)
}

func Generate(domainService string) (base.ServiceContract, error) {
	switch domainService {
	case "job":
		return new(job.Service), nil
	default:
		return nil, errors.New("Could not find domain item")
	}
}
package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/boom/app/domain/base"
	"log"
	"github.com/MoonBabyLabs/boom/app/provider"
	"strconv"
)


type Rest struct {
	*revel.Controller
}

type DataResponse struct {
	id int
}

func (c Rest) Get(domain string, resource string) revel.Result {
	data := make(map[string]interface{})
	model := base.Model{}
	model.Domain = domain
	model.Datastore = provider.Db{}.Construct()
	i, err := strconv.Atoi(resource)
	data["data"] = model.Find(i)
	data["error"] = err
	data["sucess"] = false

	return c.RenderJSON(data)
}

func (c Rest) Patch(domain string, resource string) revel.Result {
	model := base.Model{}
	item := make(map[string]interface{})
	c.Params.BindJSON(&item)
	log.Print(item)
	model.Domain = domain
	model.Datastore = provider.Db{}.Construct()
	model.Update(resource, item, true)
	data := make(map[string]interface{})
	data["success"] = true

	return c.RenderJSON(data)
}

func (c Rest) PUT(domain string, resource string) revel.Result {
	model := base.Model{}
	item := make(map[string]interface{})
	c.Params.BindJSON(&item)
	model.Domain = domain
	model.Datastore = provider.Db{}.Construct()
	model.Update(resource, item, false)
	data := make(map[string]interface{})
	data["success"] = true

	return c.RenderJSON(data)
}

func (c Rest) Post(domain string) revel.Result {
	model := base.Model{}
	item := make(map[string]interface{})
	c.Params.BindJSON(&item)
	log.Print(item)
	model.Domain = domain
	model.Datastore = provider.Db{}.Construct()
	model.Add(item)

	return c.RenderJSON(item)
}

func (c Rest) Index(domain string) revel.Result {
	model := base.Model{}
	model.Domain = domain
	model.Datastore = provider.Db{}.Construct()

	return c.RenderJSON(model.All())
}

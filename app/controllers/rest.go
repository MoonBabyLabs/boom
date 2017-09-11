package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/boom/app/domain/base"
	"log"
	"github.com/MoonBabyLabs/boom/app/provider"
	"github.com/MoonBabyLabs/boom/app/auth"
)


type Rest struct {
	*revel.Controller
}

func (c Rest) Get(domain string, resource string) revel.Result {
	cf := auth.ContentConf{}.GetContentConf(domain)

	if !cf.HasAccess(c.Request.Header, "read") {
		return c.NotFound("Unable to access resource")
	}

	model := base.Model{}
	model.Domain = domain
	model.Fields = cf.Fields
	model.Datastore = provider.Db{}.Construct()
	data := model.Find(resource)

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

func (c Rest) Put(domain string, resource string) revel.Result {
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
	cf := auth.ContentConf{}.GetContentConf(domain)

	if !cf.HasAccess(c.Request.Header,"write") {
		return c.NotFound("Unable to access page")
	}

	model := base.Model{}
	log.Print(cf.Fields)
	item := make(map[string]interface{})
	c.Params.BindJSON(&item)
	log.Print(item)
	model.Domain = domain
	model.Fields = cf.Fields
	model.Datastore = provider.Db{}.Construct()
	model.Add(item)

	return c.RenderJSON(item)
}

func (c Rest) Index(domain string) revel.Result {
	cf := auth.ContentConf{}.GetContentConf(domain)
	log.Print(c.Request.Header)

	if !cf.HasAccess(c.Request.Header,"read") {
		return c.NotFound("Unable to access page")
	}

	model := base.Model{}
	model.Domain = domain
	model.Fields = cf.Fields
	model.Datastore = provider.Db{}.Construct()

	return c.RenderJSON(model.All())
}

func (c Rest) Delete(domain string, resource string) revel.Result {
	cf := auth.ContentConf{}.GetContentConf(domain)
	data := make(map[string]interface{})

	if !cf.HasAccess(c.Request.Header, "delete") {
		return c.NotFound("unable to access page")
	}

	model := base.Model{}
	model.Domain = domain
	model.Datastore = provider.Db{}.Construct()

	if model.Delete(resource) {
		data["success"] = true

		return c.RenderJSON(data)
	}

	data["success"] = false

	return c.RenderJSON(data)
}
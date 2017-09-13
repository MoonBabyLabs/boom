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

func(c Rest) addHeaderCORS() revel.Result {
	c.Response.Out.Header().Add("Access-Control-Allow-Origin","*")

	return nil
}

func (c Rest) Get(domain string, resource string) revel.Result {
	c.addHeaderCORS()

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
	c.addHeaderCORS()

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
	c.addHeaderCORS()

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
	c.addHeaderCORS()

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
	model.Files = c.Params.Files
	model.FileManager = provider.Filemanager{}.Construct()
	model.Datastore = provider.Db{}.Construct()
	model.Add(item)

	return c.RenderJSON(item)
}

func (c Rest) Index(domain string) revel.Result {
	c.addHeaderCORS()

	cf := auth.ContentConf{}.GetContentConf(domain)

	if cf.Fields == nil {
		return c.NotFound("Either content type doesn't exist or wasnt able to parse fields correctly")
	}

	if !cf.HasAccess(c.Request.Header,"read") {
		return c.NotFound("Unable to access page")
	}

	model := base.Model{}
	model.Domain = domain
	model.Fields = cf.Fields
	model.Datastore = provider.Db{}.Construct()
	attrs := make(map[string]interface{})

	for k, v := range c.Params.Values {
		for _, b := range v {
			attrs[k] = b
		}
	}

	log.Print(attrs)
	res := model.All(attrs)

	return c.RenderJSON(res)
}

func (c Rest) Delete(domain string, resource string) revel.Result {
	c.addHeaderCORS()
	
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
// Package controllers hosts all of the API controllers available.
// Each type of API controller may have their own defined interface methods and routes.
// @todo need to refactor out access logic to middleware before getting to controller
package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/boom/app/provider"
	"github.com/MoonBabyLabs/boom/app/auth"
	"strings"
	"github.com/MoonBabyLabs/boom/app/conf"
	"github.com/MoonBabyLabs/boom/app/service/content"
	"github.com/MoonBabyLabs/kek/service"
	"strconv"
	"log"
)

type Rest struct {
	*revel.Controller
}

// Get provides the access point for GET requests to a single content resource.
// @param contentType is a string that represents the content type you would like to access.
// @param resource is the identifier for the desired resource item.
// Returns a Revel render result
func (c Rest) GetCollectionResource(contentType string, resource string) revel.Result {
	cf := auth.ContentConf{}.GetContentConf(contentType)

	if !cf.HasAccess(c.Request.Header, "read") {
		return c.NotFound("Unable to access resource")
	}

	noHistory := c.Params.Query.Get("history") == "false"
	cnt, err := provider.Content{}.Construct().Find(resource, !noHistory)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.renderContent(cnt)
}

// Options provides a route request for OPTIONS based routes.
// It is generally useful for preflight requests from browsers.
// Returns a success message in the revel render result.
func (c Rest) Options() revel.Result {
	success := make(map[string]bool)
	success["success"] = true

	return c.RenderJSON(success)
}

// Patch handles an HTTP Patch request method.
// It will patch a resource item for an individual content type if the request matches the appropriate authorization and access.
// @param contentType is a string that tells us the type of content we are patching.
// @param resource is a string that tells us the resource identifier that we are patching.
// Patch will return back a success method with data on success. It will return an appropriate HTTP error code and message on failuare.
// @todo Patch should follow some standard that also has description. Needs more research to determine how to handle standardization without overcomplicating it.
func (c Rest) PatchResource(resource string) revel.Result {
	item := make(map[string]interface{})
	c.Params.BindJSON(&item)
	log.Print(resource)
	upd, err := provider.Content{}.Construct().Update(resource, item, true)

	if err != nil {
		data := make(map[string]interface{})
		data["success"] = false
		data["error"] = err

		return c.RenderJSON(data)
	}

	return c.renderContent(upd)
}

// @todo implement
func (c Rest) PatchCollection() revel.Result {
	return c.RenderText("implement")
}

// @todo implement
func (c Rest) GetResource(resource string) revel.Result {
	revHistory := c.Params.Query.Get("_revisions") != ""
	kek, err := provider.Content{}.Construct().Find(resource, revHistory)

	if err != nil {
		return c.RenderError(err)
	}

	return c.renderContent(kek)
}

func (c Rest) PutResource(resource string) revel.Result {
	item := make(map[string]interface{})
	c.Params.BindJSON(&item)
	upd, err := provider.Content{}.Construct().Update(resource, item, false)

	if err != nil {
		data := make(map[string]interface{})
		data["success"] = false
		data["error"] = err.Error()

		return c.RenderJSON(data)
	}

	return c.renderContent(upd)
}

// @todo implement
func (c Rest) PutCollectionResource(resource string) revel.Result {
	return c.RenderText("implement")
}

// PostCollection handles an HTTP POST request from the server by creating a brand new content resource.
// The @contentType parameter tells us the type of content that we are creating.
// It returns back a json array with a succes message and the data when a new item is created.
// It will return an appropriate error code and message when the user either didn't have enough access or the system couldn't create the new content resource.
func (c Rest) PostCollectionResource(collection string) revel.Result {
	return c.RenderText("To do.. need to add collection capability for post")
}

func (c Rest) PostResource() revel.Result {
	item := make(map[string]interface{})
	c.Params.BindJSON(&item)
	kd, err := content.Default{}.Add(item, c.Params.Files)

	if err != nil {
		return c.RenderError(err)
	}

	return c.RenderJSON(kd)
}

// @todo Implement
func (c Rest) Main() revel.Result {
	attrs := make(map[string]interface{})
	limit := c.Params.Query.Get("_limit")
	history := c.Params.Query.Get("_revisions") != ""
	order := c.Params.Query.Get("_order")
	log.Print(order)
	offset, _ := strconv.Atoi(c.Params.Query.Get("_offset"))
	c.Params.Query.Del("_order")
	c.Params.Query.Del("_revisions")
	c.Params.Query.Del("_limit")
	c.Params.Query.Del("_offset")

	if limit == "" {
		limit = "20"
	}

	intLimit, _ := strconv.Atoi(limit)

	for k, v := range c.Params.Values {
		for _, b := range v {
			attrs[k] = b
		}
	}

	kekDocs, err := provider.Content{}.Construct().All(attrs, intLimit, order, offset, history)

	if err != nil {
		c.RenderError(err)
	}

	return c.RenderJSON(kekDocs)
}

func (c Rest) DeleteResource(resource string) revel.Result {
	data := make(map[string]interface{})
	err := provider.Content{}.Construct().Delete(resource)

	if err == nil {
		data["success"] = true

		return c.RenderJSON(data)
	}

	return c.NotFound(err.Error())
}

// @todo implement
func (c Rest) DeleteCollectionResource(collection string, resource string) revel.Result {
	return c.RenderText("implement")
}

// renderContent runs through a switch case to parse and return the content format filled in with the desired data.
// Parameter @cnt is a map[string]map[string]interface{} which will get parsed into the desired response.
// Param @resType is a string that tells us what type of response format that we need.
// @todo may need to refactor this out as response format types get larger and move to a more polymorphic approach instead of a switch case
// @todo need to add more response formats.
//
// Available resFormats at the moment: json, xml
func (c Rest) renderContent(cnt service.KekDoc) revel.Result {
	v := conf.Views{}
	domainPath := revel.Config.StringDefault("domain.base.path", "/")


	resType := c.Request.Header.Get("Accept")
	fmt := c.Params.Query.Get("_format")

	if fmt != "" {
		resType = fmt
	}

	if resType == "" || resType == "*/*" {
		resType = revel.Config.StringDefault("response.format.default", "application/vnd.siren+json")
	}

	if strings.Contains(resType, "json") {
		return c.RenderJSON(v.Get(resType).Run(cnt, domainPath))
	} else if strings.Contains(resType, "xml") {
		return c.RenderXML(v.Get(resType).Run(cnt, domainPath))
	} else {
		return c.Render(v.Get(resType).Run(cnt, domainPath))
	}
}
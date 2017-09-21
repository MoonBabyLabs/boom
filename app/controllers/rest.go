// Package controllers hosts all of the API controllers available.
// Each type of API controller may have their own defined interface methods and routes.
// @todo need to refactor out access logic to middleware before getting to controller
package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/boom/app/provider"
	"github.com/MoonBabyLabs/boom/app/auth"
)


type Rest struct {
	*revel.Controller
}

// Get provides the access point for GET requests to a single content resource.
// @param contentType is a string that represents the content type you would like to access.
// @param resource is the identifier for the desired resource item.
// Returns a Revel render result
func (c Rest) Get(contentType string, resource string) revel.Result {
	cf := auth.ContentConf{}.GetContentConf(contentType)

	if !cf.HasAccess(c.Request.Header, "read") {
		return c.NotFound("Unable to access resource")
	}

	noHistory := c.Params.Query.Get("history") == "false"
	cnt, err := provider.Content{}.Construct().Find(contentType, resource, !noHistory)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.renderContent(cnt, c.Params.Query.Get("_format"))
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
func (c Rest) Patch(contentType string, resource string) revel.Result {
	cf := auth.ContentConf{}.GetContentConf(contentType)

	if !cf.HasAccess(c.Request.Header, "update") {
		return c.NotFound("Unable to access page")
	}

	item := make(map[string]interface{})
	c.Params.BindJSON(&item)
	upd, err := provider.Content{}.Construct().Update(contentType, resource, item, true)
	data := make(map[string]interface{})
	data["success"] = true
	data["data"] = upd

	if err != nil {
		data["success"] = false
		data["error"] = err
	}

	return c.RenderJSON(data)
}

func (c Rest) Put(contentType string, resource string) revel.Result {
	cf := auth.ContentConf{}.GetContentConf(contentType)

	if !cf.HasAccess(c.Request.Header, "update") {
		return c.NotFound("Unable to access page")
	}

	item := make(map[string]interface{})
	c.Params.BindJSON(&item)
	upd, err := provider.Content{}.Construct().Update(contentType, resource, item, false)
	data := make(map[string]interface{})
	data["success"] = true
	data["data"] = upd

	if err != nil {
		data["success"] = false
		data["error"] = err.Error()
	}

	return c.RenderJSON(data)
}

// Post handles an HTTP POST request from the server by creating a brand new content resource.
// The @contentType parameter tells us the type of content that we are creating.
// It returns back a json array with a succes message and the data when a new item is created.
// It will return an appropriate error code and message when the user either didn't have enough access or the system couldn't create the new content resource.
func (c Rest) Post(contentType string) revel.Result {
	cf := auth.ContentConf{}.GetContentConf(contentType)

	if !cf.HasAccess(c.Request.Header,"write") {
		return c.NotFound("Unable to access page")
	}

	item := make(map[string]interface{})
	c.Params.BindJSON(&item)
	cnt, err := provider.Content{}.Construct().Add(contentType, item, c.Params.Files, cf.Fields)
	data := make(map[string]interface{})
	data["success"] = true
	data["data"] = cnt

	if err != nil {
		data["error"] = err
	}

	return c.RenderJSON(item)
}

func (c Rest) Index(contentType string) revel.Result {
	cf := auth.ContentConf{}.GetContentConf(contentType)

	if cf.Fields == nil {
		return c.NotFound("Either content type doesn't exist or wasnt able to parse fields correctly")
	}

	if !cf.HasAccess(c.Request.Header,"read") {
		return c.NotFound("Insufficient access")
	}

	attrs := make(map[string]interface{})

	for k, v := range c.Params.Values {
		for _, b := range v {
			attrs[k] = b
		}
	}

	noHistory := c.Params.Query.Get("history") == "false"
	cnt, err := provider.Content{}.Construct().All(contentType, attrs, cf.Fields, !noHistory)

	if err != nil {
		return c.NotFound(err.Error())
	}

	return c.RenderJSON(cnt)
}

func (c Rest) Delete(contentType string, resource string) revel.Result {
	cf := auth.ContentConf{}.GetContentConf(contentType)
	data := make(map[string]interface{})

	if !cf.HasAccess(c.Request.Header, "delete") {
		return c.NotFound("unable to access page")
	}
	success, err := provider.Content{}.Construct().Delete(contentType, resource)
	if success {
		data["success"] = true

		return c.RenderJSON(data)
	}

	return c.NotFound(err.Error())
}

// renderContent runs through a switch case to parse and return the content format filled in with the desired data.
// Parameter @cnt is a map[string]map[string]interface{} which will get parsed into the desired response.
// Param @resType is a string that tells us what type of response format that we need.
// @todo may need to refactor this out as response format types get larger and move to a more polymorphic approach instead of a switch case
// @todo need to add more response formats.
//
// Available resFormats at the moment: json, xml
func (c Rest) renderContent(cnt map[string]interface{}, resType string) revel.Result {
	var res revel.Result
	switch resType {
	case "json-hal":
		break
	case "xml":
		res = c.RenderXML(cnt)
		break
	case "html" :
		break
	default:
		res = c.RenderJSON(cnt)
	}

	return res
}
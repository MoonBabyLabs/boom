package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/boom/app/service/content"
	"github.com/MoonBabyLabs/kekcollections"
)

type Patch struct {
	*revel.Controller
	Base
}

// Patch handles an HTTP Patch request method.
// It will patch a resource item for an individual content type if the request matches the appropriate authorization and access.
// @param contentType is a string that tells us the type of content we are patching.
// @param resource is a string that tells us the resource identifier that we are patching.
// Patch will return back a success method with data on success. It will return an appropriate HTTP error code and message on failuare.
// TODO(miracle) Patch should follow some standard that also has description. Needs more research to determine how to handle standardization without overcomplicating it.
// TODO(miracle) Patching collections should be dynamic enough to know based on the id whether its a collection. Also, if it is a slug then we know its a collection.
func (c Patch) PatchResource(resource string) revel.Result {
	accessErr := c.HasAccess(c.Request.Header.Get("Authorization"),"update"); if accessErr != nil {
		return c.RenderError(accessErr)
	}

	item := make(map[string]interface{})
	item["id"] = resource
	c.Params.BindJSON(&item)

	if item["is_collection"] == true {
		col := kekcollections.Collection{}
		col.Id = resource
		c.Params.BindJSON(&col)
		updated, updErr := col.Patch()

		if updErr != nil {
			return c.RenderError(updErr)
		}

		return c.RenderJSON(updated)
	} else {
		upd, err := content.Default{}.Update(resource, item, true)

		if err != nil {
			data := make(map[string]interface{})
			data["success"] = false
			data["error"] = err

			return c.RenderJSON(data)
		}

		return c.RenderContent(upd)
	}
}
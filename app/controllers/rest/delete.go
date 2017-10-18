package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/boom/app/service/content"
	"github.com/MoonBabyLabs/kekcollections"
)

type Delete struct {
	*revel.Controller
	Base
}

func (c Delete) DeleteResource(resource string) revel.Result {
	accessErr := c.HasAccess(c.Request.Header.Get("jwt"),"delete"); if accessErr != nil {
		return c.RenderError(accessErr)
	}

	data := make(map[string]interface{})
	err := content.Default{}.Delete(resource)

	if err == nil {
		data["success"] = true

		return c.RenderJSON(data)
	}

	col, _ := kekcollections.Collection{}.LoadBySlug(resource, 0, false, false)
	delColErr := col.Delete(true)

	if delColErr == nil {
		data["success"] = true

		return c.RenderJSON(data)
	}

	return c.RenderError(delColErr)
}

// @todo implement
func (c Delete) DeleteCollectionResource(collection string, resource string) revel.Result {
	return c.RenderText("implement")
}

package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/boom/app/service/content"
	"github.com/MoonBabyLabs/kekcollections"
)

type Put struct {
	*revel.Controller
	Base
}

func (c Put) PutResource(resource string) revel.Result {
	accessErr := c.HasAccess(c.Request.Header.Get("Authorization"),"update"); if accessErr != nil {
		return c.RenderError(accessErr)
	}

	item := make(map[string]interface{})
	c.Params.BindJSON(&item)

	if item["is_collection"] == true {
		col := kekcollections.Collection{}
		c.Params.BindJSON(&col)
		col.Id = resource
		updatedCol, updErr := col.Replace()

		if updErr != nil {
			return c.RenderError(updErr)
		}

		return c.RenderJSON(updatedCol)
	} else {
		upd, err := content.Default{}.Update(resource, item, false)

		if err != nil {
			data := make(map[string]interface{})
			data["success"] = false
			data["error"] = err.Error()

			return c.RenderJSON(data)
		}

		return c.RenderContent(upd)
	}

}

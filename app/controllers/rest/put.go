package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/boom/app/service/content"
	"github.com/MoonBabyLabs/kekcollections"
	"github.com/MoonBabyLabs/kek"
	"errors"
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

func (c Put) PutCollectionResource(collection string, resource string) revel.Result {
	accessErr := c.HasAccess(c.Request.Header.Get("Authorization"),"update"); if accessErr != nil {
		return c.RenderError(accessErr)
	}

	kd, kdLoadErr := kek.KekDoc{}.Get(resource, false)

	if kdLoadErr != nil {
		return c.RenderError(kdLoadErr)
	}

	resType := collection[0:2]

	if resType == "dd" {
		return c.RenderError(errors.New("You are trying to save a resource into a kekdoc. This isn't possible. You must save into an appropriate kekcollection"))
	} else if resType == "cc" {
		col, colErr := kekcollections.Collection{}.LoadById(collection, false, false)

		if colErr != nil {
			return c.RenderError(colErr)
		}

		col.AddDoc(kd)
	} else {
		col, colErr := kekcollections.Collection{}.LoadBySlug(collection, 0, false, false)

		if colErr != nil {
			return c.RenderError(colErr)
		}

		col.AddDoc(kd)
	}

	return c.RenderContent(kd)
}

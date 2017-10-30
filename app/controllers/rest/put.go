package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/kekcollections"
	"github.com/MoonBabyLabs/kek"
	"errors"
	"log"
	"strings"
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
	itemType := resource[0:2]
	hiddenFieldsConf := revel.Config.StringDefault("hide.fields", "password")
	hiddenFields := strings.Split(hiddenFieldsConf, ",")

	switch itemType {
	case "cc":
		col, loadEr := kekcollections.Collection{}.LoadById(resource, false, false)

		// Looks like its not in our space. Lets try to verify that its good kekcontact and put into the space
		if loadEr != nil {

		} else {
			// Lets update the item in our space.
			c.Params.BindJSON(&col)
			col.Id = resource
			updErr := col.Save()

			if updErr != nil {
				return c.RenderError(updErr)
			}
			return c.RenderJSON(col)
		}
		break
	case "dd":
		_, loadErr := kek.Doc{}.Get(resource, false)
		log.Print(loadErr)

		if loadErr != nil {

		} else {
			upd, err := kek.Doc{}.Update(resource, item, false)

			if err != nil {
				data := make(map[string]interface{})
				data["success"] = false
				data["error"] = err.Error()

				return c.RenderJSON(data)
			}

			for _, hf := range hiddenFields {
				delete(upd.Attributes, hf)
			}

			return c.RenderContent(upd)
			break
		}
	}

	return c.RenderError(errors.New("could not find item to update"))
}

func (c Put) PutCollectionResource(collection string, resource string) revel.Result {
	accessErr := c.HasAccess(c.Request.Header.Get("Authorization"),"update"); if accessErr != nil {
		return c.RenderError(accessErr)
	}

	kd, kdLoadErr := kek.Doc{}.Get(resource, false)

	if kdLoadErr != nil {
		return c.RenderError(kdLoadErr)
	}

	resType := collection[0:2]

	if resType == "dd" {
		return c.RenderError(errors.New("you are trying to save a resource into a kekdoc. This isn't possible. You must save into an appropriate kekcollection"))
	} else if resType == "cc" {
		col, colErr := kekcollections.Collection{}.LoadById(collection, false, false)

		if colErr != nil {
			return c.RenderError(colErr)
		}

		col.AddResource(resource)
	} else {
		col, colErr := kekcollections.Collection{}.LoadBySlug(collection, false, false)

		if colErr != nil {
			return c.RenderError(colErr)
		}

		col.AddResource(resource)
	}

	return c.RenderContent(kd)
}

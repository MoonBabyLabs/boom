package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/kekcollections"
	"github.com/MoonBabyLabs/kek"
	"errors"
)

type Delete struct {
	*revel.Controller
	Base
}

func (c Delete) DeleteResource(resource string) revel.Result {
	accessErr := c.HasAccess(c.Request.Header.Get("Authorization"),"delete"); if accessErr != nil {
		return c.RenderError(accessErr)
	}

	resType := resource[0:2]
	var delErr error
	data := make(map[string]interface{})
	data["success"] = true

	switch resType {
	case "dd":
		err := kek.Doc{}.Delete(resource)

		if err == nil {
			return c.RenderJSON(data)
		} else {
			return c.RenderError(err)
		}

		break
	case "cc":
		kc := kekcollections.Collection{Id: resource}
		delErr := kc.Delete(true)

		if delErr != nil {
			return c.RenderError(delErr)
		} else {
			return c.RenderJSON(data)
		}

		break
	}

	col, loadEr := kekcollections.Collection{}.LoadBySlug(resource,false, false)

	if loadEr != nil {
		return c.RenderError(loadEr)
	}

	delErr = col.Delete(true)

	if delErr != nil {
		return c.RenderError(delErr)
	} else {
		return c.RenderJSON(data)
	}

	return c.RenderError(errors.New("Could not find " + resource + " resource to delete"))
}

func (c Delete) DeleteCollectionResource(collection string, resource string) revel.Result {
	accessErr := c.HasAccess(c.Request.Header.Get("Authorization"),"delete"); if accessErr != nil {
		return c.RenderError(accessErr)
	}

	resType := collection[0:2]

	if resType == "cc" {
		col, loadErr := kekcollections.Collection{}.LoadById(collection, false, false)

		if loadErr != nil {
			return c.RenderError(loadErr)
		}
		err := col.DeleteResource(resource)

		if err != nil {
			return c.RenderError(err)
		}

	} else {
		col, loadErr := kekcollections.Collection{}.LoadBySlug(collection, false, false)

		if loadErr != nil {
			return c.RenderError(loadErr)
		}

		err := col.DeleteResource(resource)

		if err != nil {
			return c.RenderError(err)
		}
	}

	data := make(map[string]bool)
	data["success"] = true

	return c.RenderJSON(data)
}

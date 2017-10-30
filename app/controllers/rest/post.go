package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/kekcollections"
	"github.com/MoonBabyLabs/kek"
	"errors"
	"log"
	"encoding/json"
	"strings"
)

type Post struct {
	*revel.Controller
	Base
}


// PostResource handles an HTTP POST request from the server by creating a brand new content resource.
// The @contentType parameter tells us the type of content that we are creating.
// It returns back a json array with a succes message and the data when a new item is created.
// It will return an appropriate error code and message when the user either didn't have enough access or the system couldn't create the new content resource.
func (c Post) PostResource() revel.Result {
	accessError :=c.HasAccess(c.Request.Header.Get("Authorization"),"write"); if accessError != nil {
		return c.RenderError(accessError)
	}

	item := make(map[string]interface{})
	c.Params.BindJSON(&item)

	if item["is_collection"] == true {
		marshaledCol := kekcollections.Collection{}
		json.Unmarshal(c.Params.JSON, &marshaledCol)
		savedCol, err := kekcollections.Collection{}.New(marshaledCol.Name, marshaledCol.Description, marshaledCol.ResourceIds)

		if err != nil {
			return c.RenderError(err)
		}

		return c.RenderJSON(savedCol)
	} else {
		kd, err := kek.Doc{}.New(item)

		if err != nil {
			return c.RenderError(err)
		}

		hiddenFieldsConf := revel.Config.StringDefault("hide.fields", "password")
		hiddenFields := strings.Split(hiddenFieldsConf, ",")

		for _, hf := range hiddenFields {
			delete(kd.Attributes, hf)
		}

		return c.RenderContent(kd)
	}
}


func (c Post) PostCollectionResource(collectionResource string) revel.Result {
	accessE := c.HasAccess(c.Request.Header.Get("Authorization"), "write"); if accessE != nil {
		return c.RenderError(accessE)
	}

	attrs := make(map[string]interface{})
	c.Params.BindJSON(&attrs)
	kd, newE := kek.Doc{}.New(attrs)

	if newE != nil {
		return c.RenderError(newE)
	}

	resType := collectionResource[0:2]
	col := kekcollections.Collection{}

	if resType == "dd" {
		return c.RenderError(errors.New("Can't attach posted resource to " + collectionResource + " because it is a document and not a collection"))
	} else if resType == "cc" {
		loadedCol, err := col.LoadById(collectionResource, false, false)
		log.Print(loadedCol)

		if err != nil {
			return c.RenderError(err)
		}

		loadedCol.AddResource(kd.Id)

	} else {
		// Lets assume that the collection is now a slug.
		col, err := col.LoadBySlug(collectionResource, false, false)

		if err != nil {
			return c.RenderError(err)
		}

		col.AddResource(kd.Id)
	}

	hiddenFieldsConf := revel.Config.StringDefault("hide.fields", "password")
	hiddenFields := strings.Split(hiddenFieldsConf, ",")

	for _, hf := range hiddenFields {
		delete(kd.Attributes, hf)
	}

	return c.RenderContent(kd)
}
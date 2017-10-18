package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/kekcollections"
	"github.com/MoonBabyLabs/boom/app/service/content"
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
		col := kekcollections.Collection{}
		c.Params.BindJSON(&col)
		savedCol, err := col.New()

		if err != nil {
			return c.RenderError(err)
		}

		return c.RenderJSON(savedCol)
	} else {
		item := make(map[string]interface{})
		c.Params.BindJSON(&item)
		kd, err := content.Default{}.Add(item, c.Params.Files)

		if err != nil {
			return c.RenderError(err)
		}

		return c.RenderContent(kd)
	}
}

package controllers

import (
	"github.com/revel/revel"
	"strconv"
	"github.com/MoonBabyLabs/boom/app/service/content"
	"github.com/MoonBabyLabs/kekcollections"
)

type Get struct {
	*revel.Controller
	Base
}

func (c Get) Search() revel.Result {
	if revel.Config.StringDefault("require.jwt.get", "false") == "true" {
		accessErr := c.HasAccess(c.Request.Header.Get("Authorization"),"get"); if accessErr != nil {
			return c.RenderError(accessErr)
		}
	}

	attrs := make(map[string]interface{})
	limit := c.Params.Query.Get("_limit")
	history := c.Params.Query.Get("_revisions") != ""
	order := c.Params.Query.Get("_order")
	offset, _ := strconv.Atoi(c.Params.Query.Get("_offset"))
	collections := c.Params.Query.Get("_collections")

	if collections != "" {
		cols, colLoadErr := kekcollections.Collection{}.All(true, history)

		if colLoadErr != nil {
			return c.RenderError(colLoadErr)
		}

		return c.RenderJSON(cols)
	}
	c.Params.Query.Del("_order")
	c.Params.Query.Del("_revisions")
	c.Params.Query.Del("_limit")
	c.Params.Query.Del("_offset")
	c.Params.Del("version")

	if limit == "" {
		limit = "20"
	}

	intLimit, _ := strconv.Atoi(limit)

	for k, v := range c.Params.Values {
		for _, b := range v {
			attrs[k] = b
		}
	}

	kekDocs, err := content.Default{}.All(attrs, intLimit, order, offset, history)

	if err != nil {
		c.RenderError(err)
	}

	return c.RenderJSON(kekDocs)
}


// GetResource returns a valid kekdoc, kekcollection or kekmedia item according to the resource string.
// The resource string can be the id of the resource item or a slug representing collections.
// Slugs aren't unique and you could have numerous collections having the same slug.
// @todo we need a way to differentiate and pull back the right collection resource when getting from slug. Right now, its pulling the 1st.
func (c Get) GetResource(resource string) revel.Result {

	if revel.Config.StringDefault("require.jwt.get", "false") == "true" {
		accessErr := c.HasAccess(c.Request.Header.Get("Authorization"),"get"); if accessErr != nil {
			return c.RenderError(accessErr)
		}
	}

	resType := resource[0:2]

	if resType == "dd" {
		revHistory := c.Params.Query.Get("_revisions") != ""
		kek, docErr := content.Default{}.Find(resource, revHistory)

		if docErr == nil {
			return c.RenderContent(kek)
		}
	} else if resType == "cc" {
		col, colErr := kekcollections.Collection{}.LoadById(resource, true, false)

		if colErr != nil {
			return c.NotFound("resource not found")
		}

		return c.RenderJSON(col)
	}

	col, colErr  := kekcollections.Collection{}.LoadBySlug(resource, 0, true, false)

	if colErr != nil {
		return c.NotFound("resource not found")
	}

	return c.RenderJSON(col)
}
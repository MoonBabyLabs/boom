package controllers

import (
	"github.com/revel/revel"
	"strconv"
	"github.com/MoonBabyLabs/kekcollections"
	"github.com/MoonBabyLabs/kek"
	"strings"
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
	hiddenFieldsConf := revel.Config.StringDefault("hide.fields", "password")
	hiddenFields := strings.Split(hiddenFieldsConf, ",")

	if collections != "" {
		cols, colLoadErr := kekcollections.Collection{}.All(true, history)

		if colLoadErr != nil {
			return c.RenderError(colLoadErr)
		}

		for _, col := range cols {
			for docId := range col.Docs {
				for _, hf := range hiddenFields {
					delete(col.Docs[docId].Attributes, hf)
				}
			}
		}

		return c.RenderJSON(cols)
	}
	c.Params.Del("_order")
	c.Params.Del("_revisions")
	c.Params.Del("_limit")
	c.Params.Del("_offset")
	c.Params.Del("version")
	c.Params.Del("_limit");

	if limit == "" {
		limit = "20"
	}

	intLimit, _ := strconv.Atoi(limit)

	for k, v := range c.Params.Values {
		for _, b := range v {
			attrs[k] = b
		}
	}

	q := kek.DocQuery{
		Limit: intLimit,
		Offset: offset,
		OrderBy: order,
		SearchQueries: make([]kek.SearchQuery, len(attrs)),
		WithDocRevs: history,
	}

	qCount := 0
	for k, v := range attrs {
		q.SearchQueries[qCount].Field = k
		q.SearchQueries[qCount].Value = v.(string)
		q.SearchQueries[qCount].Operator = "="
		qCount++
	}

	kekDocs, err := kek.Doc{}.Find(q)

	for kekDocId := range kekDocs {
		for _, hf := range hiddenFields {
			delete(kekDocs[kekDocId].Attributes, hf)
		}
	}

	if err != nil {
		c.RenderError(err)
	}

	return c.RenderJSON(kekDocs)
}


// GetResource returns a valid kekdoc, kekcollection or kekmedia item according to the resource string.
// The resource string can be the id of the resource item or a slug representing collections.
// Slugs aren't guaranteed unique and you could have numerous collections having the same slug.
func (c Get) GetResource(resource string) revel.Result {
	if revel.Config.StringDefault("require.jwt.get", "false") == "true" {
		accessErr := c.HasAccess(c.Request.Header.Get("Authorization"),"get"); if accessErr != nil {
			return c.RenderError(accessErr)
		}
	}

	resType := resource[0:2]
	hiddenFieldsConf := revel.Config.StringDefault("hide.fields", "password")
	hiddenFields := strings.Split(hiddenFieldsConf, ",")

	if resType == "dd" {
		revHistory := c.Params.Query.Get("_revisions") != ""
		kd, docErr := kek.Doc{}.Get(resource, revHistory)
		for _, hidden := range hiddenFields {
			delete(kd.Attributes, hidden)
		}

		if docErr == nil {
			return c.RenderContent(kd)
		}
	} else if resType == "cc" {
		col, colErr := kekcollections.Collection{}.LoadById(resource, true, false)

		if colErr != nil {
			return c.NotFound("resource not found")
		}

		for id := range col.Docs {
			for _, hf := range hiddenFields {
				delete(col.Docs[id].Attributes, hf)
			}
		}

		return c.RenderJSON(col)
	}

	col, colErr  := kekcollections.Collection{}.LoadBySlug(resource,true, false)

	if colErr != nil {
		return c.NotFound("resource not found")
	}

	for id := range col.Docs {
		for _, hf := range hiddenFields {
			delete(col.Docs[id].Attributes, hf)
		}
	}

	return c.RenderJSON(col)
}
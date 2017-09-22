package json

import (
	"strings"
	"github.com/MoonBabyLabs/boom/app/views"
)

type SirenLink struct {
	Rel []string
	Class []string
	Href string
	Title string
	Type string
}

type SirenAction struct {

}

type SirenEntity struct {
	Class []string
	Properties map[string]interface{}
	Entities []SirenEntity
	Links []SirenLink
	Actions []SirenAction
	Title string
}

type SirenResponse struct{
	Class []string
	Properties map[string]interface{}
	Links []SirenLink
	Actions []SirenAction
	Entities []SirenEntity
}

// Run returns a new view.Runner that can be used to set for specific view displays.
func (s SirenResponse) Run(
	cnt map[string]interface{},
	fields []map[string]map[string]interface{},
	urlRoute string, cType string) views.Runner {
	res := SirenResponse{}
	res.Properties = make(map[string]interface{})
	res.Class = make([]string, 1)
	res.Class[0] = cType
	res.Entities = make([]SirenEntity, 0)
	res.Actions = make([]SirenAction, 0)
	res.Properties["_cid"] = cnt["_cid"]
	res.Properties["created_at"] = cnt["created_at"]
	res.Properties["updated_at"] = cnt["updated_at"]
	res.Properties["_rev"] = cnt["_rev"]
	res.Links = make([]SirenLink, 0)
	cid, cok := cnt["_cid"].(string)

	if cok {
		slr := make([]string, 1)
		slr[0] = "self"
		selfLink := SirenLink{
			Title: cType + ": " + urlRoute + "/" + cid,
			Rel: slr,
		}
		res.Links = append(res.Links, selfLink)
	}

	if cnt["_chain"] != nil {
		chn := cnt["_chain"].([]interface{})

		for _, cv := range chn {
			for revK, data := range cv.(map[string]interface{}) {
				revi := SirenEntity{}
				revi.Links = make([]SirenLink, 1)
				li := SirenLink{}
				li.Rel = make([]string, 1)
				li.Title = "Revision: " + revK
				dm, ok := data.(map[string]interface{})

				if ok && dm != nil {
					revi.Properties = dm
				}

				if cnt["_cid"] != nil {
					li.Href = urlRoute + cnt["_cid"].(string) + "?_rev=" + revK
				}
				li.Rel[0] = "self"
				revi.Links = append(revi.Links, li)
				res.Entities = append(res.Entities, revi)
				revi.Title = "Revision: " + revK
			}
		}
	}

	for _, v := range fields {
		for f, d := range v {
			if cnt[f] == nil {
				continue
			}

			cntType := d["type"]

			if cntType == nil {
				res.Properties[f] = cnt[f]
			}

			ft := cntType .(string)

			if strings.Contains(ft, "[]link") {
				al := cnt[f].([]map[string]string)
				for _, l := range al {
					newLink := SirenLink{}
					newLink.Href = l["href"]
					newLink.Title = l["title"]
					newLink.Rel = append(newLink.Rel, "about")
					res.Links = append(res.Links, newLink)
				}

				continue
			} else if strings.Contains(ft, "link") {
				mpf := cnt[f].(map[string]string)
				newLink := SirenLink{}
				newLink.Href = mpf["href"]
				newLink.Title = mpf["title"]
				newLink.Rel = append(newLink.Rel, "about")
				res.Links = append(res.Links, newLink)
			} else if strings.Contains(ft, "relation") {

			} else {
				res.Properties[f] = cnt[f]
			}
		}
	}

	return res
}
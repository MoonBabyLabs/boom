package json

import (
	"strings"
	"github.com/MoonBabyLabs/boom/app/views"
	"log"
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

	// Lets add the revision chain
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

	// Lets add the properties and media
	for _, v := range fields {
		for f, d := range v {
			if cnt[f] == nil {
				continue
			}

			cntType := d["type"]

			if cntType == nil {
				res.Properties[f] = cnt[f]
			}

			ft := cntType.(string)
			log.Print(ft)

			if strings.Contains(ft, "link") {
				log.Print(cnt[f])
				al, isSlice := cnt[f].([]interface{})

				// Lets handle for non-sliced links
				if !isSlice {

					mpf, isMap := cnt[f].(map[string]interface{})

					if isMap {
						log.Panic("Cannot resolve this link structure for siren:", cnt[f])
					}

					newLink := SirenLink{}
					log.Print(mpf)
					newLink.Href, _ = mpf["href"].(string)
					newLink.Title, _ = mpf["title"].(string)
					newLink.Rel = make([]string, 1)
					newLink.Rel[0] = "about"
					res.Links = append(res.Links, newLink)
				} else {
					for _, l := range al {
						nl := l.(map[string]interface{})
						newLink := SirenLink{}
						newLink.Href, _ = nl["href"].(string)
						newLink.Title, _ = nl["title"].(string)
						newLink.Rel = make([]string, 1)

						if nl["rel"] != nil {
							newLink.Rel[0] = nl["rel"].(string)
						} else {
							newLink.Rel[0] = "about"
						}

						res.Links = append(res.Links, newLink)
					}
				}
			} else if strings.Contains(ft, "relation") {

			} else {
				res.Properties[f] = cnt[f]
			}
		}
	}

	return res
}
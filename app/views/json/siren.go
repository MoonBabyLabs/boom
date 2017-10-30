package json

import (
	"github.com/MoonBabyLabs/boom/app/views"
	"github.com/MoonBabyLabs/kek"
)

type SirenLink struct {
	Rel []string `json:"rel"`
	Class []string `json:"class"`
	Href string `json:"href"`
	Title string `json:"title"`
	Type string `json:"type"`
}

type SirenAction struct {
	Name string `json:"name"`
	Class []string `json:"class"`
	Method string `json:"method"`
	Href string `json:"href"`
	Title string `json:"title"`
	Type string `json:"type"`
	Fields []map[string]string `json:"fields"`
}

type SirenEntity struct {
	Class []string `json:"class"`
	Properties map[string]interface{} `json:"properties"`
	Entities []SirenEntity `json:"entities"`
	Links []SirenLink  `json:"links"`
	Actions []SirenAction `json:"actions"`
	Title string `json:"title"`
}

type SirenResponse struct{
	Class []string `json:"class"`
	Properties map[string]interface{} `json:"properties"`
	Links []SirenLink `json:"links"`
	Actions []SirenAction `json:"actions"`
	Entities []SirenEntity `json:"entities"`
}

// Run returns a new view.Runner that can be used to set for specific view displays.
func (s SirenResponse) Run(cnt kek.Doc, urlRoute string) views.Runner {
	res := SirenResponse{}
	res.Properties = make(map[string]interface{})
	res.Entities = make([]SirenEntity, 0)
	res.Actions = []SirenAction{
		{Name:"Put Resource", Method:"PUT", Href:"/" + cnt.Id, Class: []string{}, Fields: []map[string]string{}},
		{Name:"Patch Resource", Method:"PATCH", Href:"/" + cnt.Id, Class: []string{}, Fields: []map[string]string{}},
		{Name:"Delete Resource", Method:"DELETE", Href:"/" + cnt.Id, Class: []string{}, Fields: []map[string]string{}},
	}
	res.Properties["_kid"] = cnt.Id
	res.Properties["created_at"] = cnt.CreatedAt
	res.Properties["updated_at"] = cnt.UpdatedAt
	res.Properties["_rev"] = cnt.Rev
	res.Links = make([]SirenLink, 0)
	links := make(map[string]SirenLink)
	selfLink := SirenLink{
		Title: cnt.Id,
		Rel: []string{"self"},
		Href: urlRoute + "/" + cnt.Id,
	}
	res.Links = append(res.Links, selfLink)

	// Lets add the revision chain
	if cnt.Revisions != nil {
		for _, block := range cnt.Revisions.GetBlocks() {
			revi := SirenEntity{}
			revi.Class = []string{"revision"}
			revi.Actions = []SirenAction{}
			revi.Links = []SirenLink{
				{
					Rel: []string{"self"},
					Title: "Revision: " + block.HashString(),
					Href: "/" + cnt.Id + "?_rev=" + block.HashString(),
				},
			}
			revi.Title = "Revision: " + block.HashString()
			res.Entities = append(res.Entities, revi)
		}
	}

	// Lets add the properties and media
	for f, v := range cnt.Attributes {
		vMap, isMap := v.(map[string]string)

		if isMap && vMap["type"] == "link" {
			newLink := SirenLink{}
			newLink.Href = vMap["href"]
			newLink.Class = make([]string, 1)
			newLink.Class[0] = vMap["class"]
			newLink.Rel = make([]string, 1)
			newLink.Rel[0] = vMap["relationship"]
			newLink.Title = vMap["title"]
			links[f] = newLink
		} else {
			res.Properties[f] = v
		}
	}

	for _, link := range links {
		res.Links = append(res.Links, link)
	}

	return res
}
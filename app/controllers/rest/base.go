package controllers

import (
	"github.com/MoonBabyLabs/boom/app/conf"
	"strings"
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/kek"
	"github.com/MoonBabyLabs/kekaccess"
)

type Base struct {
	*revel.Controller
}


// renderContent runs through a switch case to parse and return the content format filled in with the desired data.
// Parameter @cnt is a map[string]map[string]interface{} which will get parsed into the desired response.
// Param @resType is a string that tells us what type of response format that we need.
// @todo may need to refactor this out as response format types get larger and move to a more polymorphic approach instead of a switch case
// @todo need to add more response formats.
//
// Available resFormats at the moment: json, xml
func (c Base) RenderContent(cnt kek.KekDoc) revel.Result {
	v := conf.Views{}
	domainPath := revel.Config.StringDefault("domain.base.path", "/")


	resType := c.Request.Header.Get("Accept")
	fmt := c.Params.Query.Get("_format")

	if fmt != "" {
		resType = fmt
	}

	if resType == "" || resType == "*/*" {
		resType = revel.Config.StringDefault("response.format.default", "application/vnd.siren+json")
	}

	if strings.Contains(resType, "json") {
		return c.RenderJSON(v.Get(resType).Run(cnt, domainPath))
	} else if strings.Contains(resType, "xml") {
		return c.RenderXML(v.Get(resType).Run(cnt, domainPath))
	} else {
		return c.Render(v.Get(resType).Run(cnt, domainPath))
	}
}

func (c Base) HasAccess(jToken string, actionMethod string) error {
	jToken = strings.Replace(jToken, "Bearer ", "", 1)
	return kekaccess.Access{}.ValidateJwtAccess(jToken, actionMethod)
}

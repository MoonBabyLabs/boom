package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/kekcontact"
	"encoding/json"
	"github.com/MoonBabyLabs/kekaccess"
	"github.com/MoonBabyLabs/kekspace"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Token() revel.Result {
	content := make(map[string]string)
	json.Unmarshal(c.Params.JSON, &content)
	passedSecret := content["access_token"]
	delete(content, "access_token")

	token, err := kekaccess.Access{}.NewJwt(passedSecret, content)

	if err != nil {
		c.RenderError(err)
	}

	jsonString := map[string]string{}
	jsonString["token"] = token

	return c.RenderJSON(jsonString)
}

func (c App) Install() revel.Result {
	ks, _ := kekspace.Kekspace{}.Load()

	if ks.Name != "" {
		return c.NotFound("page not found")
	}

	return c.Render()
}

func (c App) SaveInstall() revel.Result {
	data := c.Params.Form
	owner := kekcontact.Contact{}
	owner.Name = data.Get("name")
	owner.Email = data.Get("email")
	owner.Phone = data.Get("phone")
	owner.Address = data.Get("address")
	owner.City = data.Get("city")
	owner.CountryCode = data.Get("country_code")
	owner.PostalCode = data.Get("postal_code")
	owner.Region = data.Get("region")
	owner.Company = kekcontact.Company{
		Name: data.Get("company_name"),
		Phone: data.Get("company_phone"),
		Email: data.Get("company_email"),
		Address: data.Get("company_address"),
		City: data.Get("company_city"),
		CountryCode: data.Get("company_country_code"),
		PostalCode: data.Get("company_postal_code"),
		Region: data.Get("company_region"),
	}
	ks, ksErr := kekspace.Kekspace{}.New(data.Get("kekspace"), "", owner, []kekcontact.Contact{owner})

	if ksErr != nil {
		return c.RenderError(ksErr)
	}

	ka := kekaccess.Access{}
	ka.GenerateSecret(ks.KekId + "/")
	token := ka.AddAccessToken(true, true, true, true, true)

	return c.RenderHTML("<p>Your KekBoom Installation was a success. Store this Access code in a safe place. Your apps will need to pass it through via different requests: <br />" + token + "</p>")
}

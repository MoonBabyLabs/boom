package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/kek/service"
	"io/ioutil"
	"encoding/json"
	service2 "github.com/MoonBabyLabs/boom/app/service"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Token() revel.Result {
	token, err := service2.JWT{}.New()
	secret, err := service2.JWT{}.GetSecret()

	if err != nil {
		c.RenderError(err)
	}

	tokenString, tErr := token.SignedString(secret)

	if tErr != nil {
		return c.RenderError(tErr)
	}
	jsonString := map[string]string{}
	jsonString["token"] = tokenString

	return c.RenderJSON(jsonString)
}

func (c App) Install() revel.Result {
	ks := service.Kekspace{}
	service.Load(service.KEK_SPACE_CONFIG, &ks)
	sample := make(map[string]interface{})

	if ks.Name != "" {
		return c.NotFound("page not found")
	}

	confFile, confErr := ioutil.ReadFile("conf/kekspace.json")

	if confErr != nil {
		return c.RenderError(confErr)
	}

	json.Unmarshal(confFile, &ks)
	json.Unmarshal(confFile, &sample)
	ksErr := service.Save(service.KEK_SPACE_CONFIG, ks)

	if ksErr != nil {
		return c.RenderError(ksErr)
	}

	service2.JWT{}.GenerateSecret()

	return c.RenderText("success")
}

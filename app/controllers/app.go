package controllers

import (
	"github.com/revel/revel"
	"github.com/MoonBabyLabs/kek/service"
	"io/ioutil"
	"go/build"
	"encoding/json"
	"github.com/satori/go.uuid"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Install() revel.Result {
	ks := service.Kekspace{}
	ksConf := service.KekspaceConfig{}
	loaded, _ := ks.Load()

	if loaded.Id != uuid.Nil {
		return c.NotFound("page not found")
	}

	confFile, confErr := ioutil.ReadFile(build.Default.GOPATH + "/src/github.com/MoonBabyLabs/boom/conf/kekspace.json")

	if confErr != nil {
		return c.RenderError(confErr)
	}

	json.Unmarshal(confFile, ksConf)

	_, ksErr := service.Kekspace{}.New(ksConf)

	if ksErr != nil {
		return c.RenderError(ksErr)
	}

	return c.RenderText("success")
}

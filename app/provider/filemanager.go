package provider

import (
	"github.com/MoonBabyLabs/boom/app/service/filemanager"
	"github.com/revel/revel"
)

type Filemanager struct {

}


func (fm Filemanager) Construct() filemanager.Contract {
	manager := filemanager.Local{}
	manager.Path = revel.Config.StringDefault("filemanager.path", "/files/")

	return manager
}
package app

import (
	"github.com/MoonBabyLabs/boom/app/views"
	"github.com/MoonBabyLabs/boom/app/views/json"
	"log"
)

type Config struct {
}


type Views struct {
	List map[string]views.Runner
}

func (v Views) SetList() Views {
	v.List = make(map[string]views.Runner)

	v.List["application/json"] = json.SirenResponse{}
	v.List["application/vnd.siren+json"] = json.SirenResponse{}

	// Modify but do not delete the default option or you will have problems
	v.List["default"] = json.SirenResponse{}

	return v
}

func (v Views) Get(vtype string) views.Runner {
	v = v.SetList()
	log.Print(vtype)

	if v.List[vtype] == nil {
		return v.List["default"]
	}

	return v.List[vtype]
}
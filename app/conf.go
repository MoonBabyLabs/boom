package app

import (
	"github.com/MoonBabyLabs/boom/app/views"
	"github.com/MoonBabyLabs/boom/app/views/json"
	"log"
	"github.com/MoonBabyLabs/boom/app/service/publisher"
)

type Config struct {
}


type Views struct {
	List map[string]views.Runner
}

type Publishers struct {
	List map[string]publisher.Manager
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

func (p Publishers) GetPublishers(publishers []string) []publisher.Manager {
	list := make([]publisher.Manager, len(publishers))

	for k, pub := range publishers {
		list[k] = Publishers{}.Find(pub)
	}

	return list
}

func (pi Publishers) Find(p string) publisher.Manager {
	up := FindPublisher(p)

	return up
}
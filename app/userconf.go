package app

import (
	"github.com/MoonBabyLabs/boom/app/service/publisher"
	"github.com/MoonBabyLabs/boom/app/views"
	"github.com/MoonBabyLabs/boom/app/views/json"
)

func FindPublisher(p string) publisher.Manager {
	list := make(map[string]publisher.Manager)


	return list[p]
}

func FindView(v string) views.Runner {
	list := make(map[string]views.Runner)
	list["json"] = json.SirenResponse{}

	return list[v]
}
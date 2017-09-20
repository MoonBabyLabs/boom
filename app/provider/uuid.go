package provider

import (
	"github.com/MoonBabyLabs/boom/app/service/uuid"
)

type Uuid struct {

}

func (u Uuid) Construct() uuid.Generator {
	return uuid.Default{}
}

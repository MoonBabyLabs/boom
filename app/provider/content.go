package provider

import (
	"github.com/MoonBabyLabs/boom/app/service/content"
)

type Content struct {

}

func (c Content) Construct() content.Manager {
	return content.Default{}
}

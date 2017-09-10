package service

import (
	"github.com/MoonBabyLabs/boom/app/datastore"
	"github.com/MoonBabyLabs/boom/app/domain/base"
)

type Content struct {
	Model base.ModelContract
	Db datastore.Contract
}

func (c *Content) Init() Contract {

	return c
}
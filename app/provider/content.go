package provider

import (
	"reflect"
	"github.com/MoonBabyLabs/boom/app/service"
)

type Content struct {

}


func (c Content) Construct() service.Contract {
	content := service.Content{}

	return content
}


func (db Content) GetName() string {
	return reflect.TypeOf(db).Name()
}

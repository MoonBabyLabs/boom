package provider

import (
	"gigdub/app/datastore"
	"reflect"
	"github.com/revel/revel"
)

type Db struct {
}

func (db Db) Construct() datastore.Contract {
	ds := datastore.Mongo{}
	ds.Init(revel.Config.StringDefault("db.name", "content"), revel.Config.StringDefault("db.host", "localhost:27017"))

	return ds
}

func (db Db) GetName() string {
	return reflect.TypeOf(db).Name()
}

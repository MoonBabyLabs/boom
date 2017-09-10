package provider

import (
	"github.com/MoonBabyLabs/boom/app/datastore"
	"reflect"
	"github.com/revel/revel"
)

type Db struct {
}

func (db Db) Construct() datastore.Contract {
	ds := datastore.Tiedot{}
	ds.DbName = revel.Config.StringDefault("db.name", "content")
	ds.DbHost = revel.Config.StringDefault("db.host", "/tmpDb")
	return ds.Init(ds.DbName, ds.DbHost)
}

func (db Db) GetName() string {
	return reflect.TypeOf(db).Name()
}

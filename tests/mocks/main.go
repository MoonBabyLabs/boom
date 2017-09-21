package mocks

import "github.com/MoonBabyLabs/boom/app/datastore"

type Datastore struct {

}

func (d *Datastore) Init(dbName string, dbConnection string) datastore.Contract {
	return d
}

func (d *Datastore) Connect(a string, b string) datastore.Contract {
	return d
}

func (d *Datastore) Find(c string, r string) map[string]interface{} {
	de := make(map[string]interface{})
	de["sample"] = "content"

	return de
}

func (d *Datastore) Insert(c string, r map[string]interface{}) bool {
	return true
}

func (d *Datastore) SetDomain(domain string) datastore.Contract {
	return d
}

func (d *Datastore) GetDomain() string {
	return "test"
}

func (d *Datastore) All(collection string, query datastore.Query) []map[string]interface{} {
	s := make([]map[string]interface{}, 0)
	f := make(map[string]interface{})
	f["sample"] = "content"
	s = append(s, f)

	return s
}

func (d *Datastore) Update(colleciton string, resource string, content map[string]interface{}, patch bool) bool {
	return true
}

func (d *Datastore) Delete(a string, r string) bool {
	return true
}
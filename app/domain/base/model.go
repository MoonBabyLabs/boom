package base

import "gigdubservice/app/datastore"

type ModelContract interface {
}

type DataStoreContract interface {
	Init(dbName string, dbConnection string, domain string) DataStoreContract
	Connect(connection string, dbName string) DataStoreContract
	Find(model Model) Model
	Insert(model Model) Model
}

type Model struct {
	Datastore DataStoreContract
	ResourceFinder datastore.ResourceFinder
	Domain string
	Entities []struct{}
}

func (m Model) init(datastore DataStoreContract, resourceFinder datastore.ResourceFinder) Model {
	active := Model{Datastore:datastore, ResourceFinder:resourceFinder}
	m.Datastore = datastore.Init()

	return active
}

func (m Model) find() []struct{} {
	m.Datastore.Find(m)

	return m.Entities
}
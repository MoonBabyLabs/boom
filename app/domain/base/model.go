package base

import "gigdub/app/datastore"

type ModelContract interface {
}

type EntityContract interface {

}

type Entities map[string] interface{}

type Model struct {
	Datastore datastore.Contract
	ResourceFinder datastore.ResourceFinder
	Domain string
	Entities Entities
}

func (m Model) init(datastore datastore.Contract, resourceFinder datastore.ResourceFinder) Model {
	active := Model{Datastore:datastore, ResourceFinder:resourceFinder}
	m.Datastore = datastore

	return active
}

func (m Model) Find(resource interface{}) Entities {
	m.Datastore.Find(resource)

	return m.Entities
}

func (m Model) GetDomain() string {
	return m.Domain
}

func (m Model) SetDomain(domain string) Model {
	m.Domain = domain

	return m
}

func (m Model) Add(items ...interface{}) bool {
	if m.Datastore.Insert(items) {
		return true
	}

	return false
}
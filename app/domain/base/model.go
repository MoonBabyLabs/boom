package base

import (
	"github.com/MoonBabyLabs/boom/app/datastore"
	"log"
	"time"
)

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
	Fields []map[string]map[string]interface{}
}

func (m Model) init(datastore datastore.Contract, resourceFinder datastore.ResourceFinder) Model {
	active := Model{Datastore:datastore, ResourceFinder:resourceFinder}
	m.Datastore = datastore

	return active
}

func (m Model) Find(resource string) map[string]interface{} {
	return m.Datastore.Find(m.Domain, resource)
}

func (m Model) GetDomain() string {
	return m.Domain
}

func (m Model) SetDomain(domain string) Model {
	m.Domain = domain

	return m
}

func (m Model) Add(items map[string]interface{}) bool {
	entity := make(map[string]interface{})
	entity["updated_at"] = time.Now()
	entity["created_at"] = time.Now()

	for k, v := range m.Fields {
		log.Print(k)
		for p, n := range v {
			if items[p] != nil {
				entity[p] = items[p]
			}
			// Can be used for later
			log.Print(n)
		}
	}
	if m.Datastore.Insert(m.Domain, entity) {
		return true
	}

	return false
}

func (m Model) Delete(resource string) bool {
	return m.Datastore.Delete(m.Domain, resource)
}

func (m Model) All() []map[string]interface{} {
	return m.Datastore.All(m.Domain)
}

func (m Model) Update(resource string, content map[string]interface{}, patch bool) Model {
	m.Datastore.Update(m.Domain, resource, content, patch)

	return m
}
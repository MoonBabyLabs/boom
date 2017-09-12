package base

import (
	"github.com/MoonBabyLabs/boom/app/datastore"
	"time"
	"github.com/MoonBabyLabs/boom/app/service/filemanager"
	"mime/multipart"
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
	Files map[string][]*multipart.FileHeader
	FileManager filemanager.Contract
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
	entity := m.FileManager.Add(m.Files, m.Fields)
	entity["updated_at"] = time.Now()
	entity["created_at"] = time.Now()

	for _, v := range m.Fields {
		for p, _ := range v {
			if items[p] != nil && entity[p] == nil {
				entity[p] = items[p]
			}
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

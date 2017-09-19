package base

import (
	"github.com/MoonBabyLabs/boom/app/datastore"
	"time"
	"github.com/MoonBabyLabs/boom/app/service/filemanager"
	"mime/multipart"
	"github.com/MoonBabyLabs/boom/app/service/chain"
	"encoding/json"
	"errors"
	"log"
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
	Entity map[string]interface{}
	Fields []map[string]map[string]interface{}
	Chain chain.BoomChainHandler
}

func (m Model) init(datastore datastore.Contract, resourceFinder datastore.ResourceFinder) Model {
	active := Model{Datastore:datastore, ResourceFinder:resourceFinder}
	m.Datastore = datastore

	return active
}

func (m Model) Find(resource string, revHistory bool) map[string]interface{} {
	m.Entity = m.Datastore.Find(m.Domain, resource)
	m.Entity["_rev"] = ""

	if m.Entity["_chain"] != nil && !revHistory {
		delete(m.Entity, "_chain")
	}

	return m.Entity
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
	baseEnt, _ := json.Marshal(entity)
	block := chain.BoomBlock{}.SetTimestamp(time.Now().Unix()).SetIndex(0).SetData(baseEnt).SetAuthor("", "", "asdf3333asdf")
	entity["_chain"] = m.Chain.AddBlock(block).Blocks()
	entity["_rev"] = m.Chain.Block().HashString()

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

func (m Model) All(attributes map[string]interface{}) []map[string]interface{} {
	q := datastore.Query{}
	where := datastore.WhereQuery{}

	for _, v := range m.Fields {
		for f, _ := range v {
			if attributes[f] != nil {
				where.Reference = f
				where.Value = attributes[f]
			}
		}
	}

	if where.Value != nil {
		q.Where = where
	}

	return m.Datastore.All(m.Domain, q)
}

func (m Model) Update(resource string, content map[string]interface{}, patch bool) (Model, error) {
	item := m.Datastore.Find(m.Domain, resource)

	if len(item) == 0 {
		return m, errors.New("Could find resource. Please confirm you have the correct key")
	}

	if item["_rev"] != content["_rev"] {
		return m, errors.New("_rev key does not match for last updated key. Please pull latest content and try again")
	}

	content["updated_at"] = time.Now()
	m.Entity = content
	log.Print(content["_rev"])
	rev, _ := content["_rev"].(string)
	log.Print(item["_chain"])
	cn, _:= item["_chain"].([]interface{})
	newCn, err, hash := m.updateChain(cn, rev)
	log.Print(hash)

	if err != nil {
		return m, err
	}

	content["_chain"] = newCn
	content["_rev"] = hash
	log.Print(content["_rev"])
	m.Datastore.Update(m.Domain, resource, content, patch)

	delete(m.Entity, "_chain")

	return m, nil
}

func (m Model) updateChain(cn []interface{}, lastRev string) ([]interface{}, error, string) {

	index := len(cn)
	log.Print(index)
	log.Print(cn[index - 1])
	log.Print(cn[10])
	lastRevBlock := cn[index - 1].(map[string]interface{})
	var lastDocdRev string

	for hash, _ := range lastRevBlock {
		lastDocdRev = hash
	}

	log.Print(lastDocdRev)
	log.Print(lastRev)

	if lastDocdRev != lastRev {
		err := errors.New("_rev key does not match the last updated key. Please update your content to the latest instance and try again")

		return nil, err, ""
	}

	data, _ := json.Marshal(m.Entity)
	b := chain.BoomBlock{}.SetData(data).SetAuthor("", "", "asl3933wa").SetTimestamp(time.Now().Unix()).SetIndex(index)
	b.SetPreviousHash(lastDocdRev)
	rev := make(map[string]interface{})
	rev[b.HashString()] = data
	cn = append(cn, rev)

	return cn, nil, b.HashString()
}

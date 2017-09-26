package content

import (
	"time"
	"encoding/json"
	"github.com/MoonBabyLabs/boom/app/service/chain"
	"errors"
	"log"
	"github.com/MoonBabyLabs/boom/app/datastore"
	"github.com/MoonBabyLabs/boom/app/service/filemanager"
	"mime/multipart"
	"github.com/MoonBabyLabs/boom/app/service/uuid"
	"github.com/MoonBabyLabs/boom/app/service/publisher"
)

type Default struct {
	chain chain.BoomChainHandler
	store datastore.Contract
	fields []map[string]map[string]interface{}
	publisher []publisher.Manager
	fileManager filemanager.Contract
	domain string
	uuidGenerator uuid.Generator
}

func (m Default) Find(contentType string, resource string, revHistory bool) (map[string]interface{}, error) {
	ent := m.Datastore().Find(contentType, resource, revHistory)

	if ent["_chain"] != nil && !revHistory {
		delete(ent, "_chain")
	}

	return ent, nil
}

func (m Default) Add(
	contentType string,
	items map[string]interface{},
	files map[string][]*multipart.FileHeader,
	fields []map[string]map[string]interface{}) (map[string]interface{}, error) {

	entity := m.FileManager().Add(files, fields)
	entity["updated_at"] = time.Now()
	entity["created_at"] = time.Now()
	baseEnt, _ := json.Marshal(entity)
	block := chain.BoomBlock{}.
		SetTimestamp(time.Now().Unix()).
		SetIndex(0).SetData(baseEnt).
		SetAuthor("", "", "asdf3333asdf")
	entity["_chain"] = m.Chain().AddBlock(block).Blocks()
	entity["_rev"] = block.HashString()
	entity["_cid"] = m.UuidGenerator().New(block.HashString())

	for _, v := range fields {
		for p := range v {
			if items[p] != nil && entity[p] == nil {
				entity[p] = items[p]
			}
		}
	}

	if m.Datastore().Insert(contentType, entity) {
		for _, p := range m.publisher {
			p.Insert(contentType, entity)
		}

		return entity, nil
	}

	return entity, errors.New("Failed to add new content item")
}

// Delete removes an instance of a content type
func (m Default) Delete(contentType string, resource string) (bool, error) {
	success := m.Datastore().Delete(contentType, resource)

	if success {
		for _, p :=range m.publisher {
			p.Delete(contentType, resource)
		}

		return success, nil
	}

	return false, errors.New("Could not delete record from database")
}

// All retrieves content documents that match the requested criteria.
func (m Default) All(
	contentType string,
	attributes map[string]interface{},
	fields []map[string]map[string]interface{},
	revHistory bool) ([]map[string]interface{}, error) {
	q := datastore.Query{}
	where := datastore.WhereQuery{}

	for _, v := range fields {
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

	cnt := m.Datastore().All(contentType, q, revHistory)

	if len(cnt) > 0 {
		return cnt, nil
	}

	return cnt, errors.New("Failed to find any content resources")
}

// Update sends the content to the datastore to save a single content item. It will also update the necessary chains for that content piece.
func (m Default) Update(
	contentType string,
	resource string,
	content map[string]interface{},
	patch bool) (map[string]interface{}, error) {

	item := m.Datastore().Find(contentType, resource, true)

	if len(item) == 0 {
		return content, errors.New("Could not find resource. Please confirm you have the correct key")
	}

	log.Print(item["_rev"])
	log.Print(content["_rev"])

	if item["_rev"] != content["_rev"] {
		return content, errors.New("_rev key does not match for last updated key. Please pull latest content and try again")
	}
	content["updated_at"] = time.Now()
	rev, _ := content["_rev"].(string)
	cn, _:= item["_chain"].([]interface{})
	newChn, err, hash := m.updateChain(cn, rev, content)
	if err != nil {
		return content, err
	}
	content["_chain"] = newChn
	content["_rev"] = hash
	m.Datastore().Update(contentType, resource, content, patch)

	for _, p := range m.publisher {
		p.Update(contentType, content, patch)
	}

	return content, nil
}

// SetChain sets the instantiated ChainHandler for the content manager.
func (m Default) SetChain(handler chain.BoomChainHandler) Manager {
	m.chain = handler

	return m
}

// Chain returns back instance of ChainHandler being used.
func (m Default) Chain() chain.BoomChainHandler {
	return m.chain
}

// Datastore returns back the instance of the Datastorer being used.
func (m Default) Datastore() datastore.Contract {
	return m.store
}

// SetDatastore sets the default backend data storage for the Boom App and returns instance of the Content Manager.
func (m Default) SetDatastore(store datastore.Contract) Manager {
	m.store = store

	return m
}

// FileManager returns the default file manager for the Boom App and returns instance of the Content Manager.
func (m Default) FileManager() filemanager.Contract {
	return m.fileManager
}

// SetFileManager sets the default file manager for the Boom App.
func (m Default) SetFileManager(fileManager filemanager.Contract) Manager {
	m.fileManager = fileManager

	return m
}

func (m Default) SetUuidGenerator(generator uuid.Generator) Manager {
	m.uuidGenerator = generator

	return m
}

func (m Default) SetPublishers(publishers []publisher.Manager) Manager {
	m.publisher = publishers

	return m
}

func (m Default) UuidGenerator() uuid.Generator {
	return m.uuidGenerator
}

// updateChain is a convenience function to help with updating a chain.
func (m Default) updateChain(
	cn []interface{},
	lastRev string,
	entity map[string]interface{}) ([]interface{}, error, string) {

	index := len(cn)
	data, _ := json.Marshal(entity)
	b := m.Chain().Block().SetData(data).SetAuthor("", "", "asl3933wa").SetTimestamp(time.Now().Unix()).SetIndex(index)
	b.SetPreviousHash(entity["_rev"].(string))
	rev := make(map[string]interface{})
	rev[b.HashString()] = data
	cn = append(cn, rev)

	return cn, nil, b.HashString()
}
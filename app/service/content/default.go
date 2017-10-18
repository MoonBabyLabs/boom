package content

import (
	"github.com/MoonBabyLabs/boom/app/service/filemanager"
	"mime/multipart"
	"github.com/MoonBabyLabs/boom/app/service/publisher"
	"github.com/MoonBabyLabs/kek"
)

type Default struct {
	fields []map[string]map[string]interface{}
	publisher []publisher.Manager
	fileManager filemanager.Contract
}

func (m Default) Find(resource string, revHistory bool) (kek.KekDoc, error) {
	kd := kek.KekDoc{}
	item, getErr := kd.Get(resource, revHistory)

	return item, getErr
}

func (m Default) Add(
	attrs map[string]interface{},
	files map[string][]*multipart.FileHeader) (kek.KekDoc, error) {

	return kek.KekDoc{}.New(attrs)
}

// Delete removes an instance of a content type
func (m Default) Delete(resource string) error {
	return kek.KekDoc{}.Delete(resource)
}

// All retrieves content documents that match the requested criteria.
func (m Default) All(attributes map[string]interface{}, limit int, order string, offset int, revHistory bool) ([]kek.KekDoc, error) {
		dq := kek.DocQuery{}
		dq.SearchQueries = make([]kek.SearchQuery, len(attributes))
		dq.OrderBy = order
		dq.Offset = offset
		count := 0
		dq.Limit = limit


		for field, val := range attributes {
			dq.SearchQueries[count].Value = val.(string)
			dq.SearchQueries[count].Field = field
			dq.SearchQueries[count].Operator = "="
			count++
		}

		return kek.KekDoc{}.Find(dq)
}

// Update sends the content to the datastore to save a single content item. It will also update the necessary chains for that content piece.
func (m Default) Update(
	resource string,
	attrs map[string]interface{},
	patch bool) (kek.KekDoc, error) {
		item := kek.KekDoc{}

		return item.Update(resource, attrs, patch)
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
package content

import (
	"github.com/MoonBabyLabs/boom/app/service/chain"
	"github.com/MoonBabyLabs/boom/app/datastore"
	"github.com/MoonBabyLabs/boom/app/service/filemanager"
	"mime/multipart"
	"github.com/MoonBabyLabs/boom/app/service/uuid"
)


type Manager interface {
	Find(
		contentType string,
		resource string,
		revHistory bool) (map[string]interface{}, error)

	Add(
		contentType string,
		items map[string]interface{},
		files map[string][]*multipart.FileHeader,
		fields []map[string]map[string]interface{}) (map[string]interface{}, error)

	Delete(contentType, resource string) (bool, error)

	All(
		contentType string,
		attributes map[string]interface{},
		fields []map[string]map[string]interface{},
		revHistory bool) ([]map[string]interface{}, error)

	Update(
		group string,
		resource string,
		content map[string]interface{},
		patch bool) (map[string]interface{}, error)

	Chain() chain.BoomChainHandler

	SetChain(handler chain.BoomChainHandler) Manager

	SetDatastore(store datastore.Contract) Manager

	Datastore() datastore.Contract

	FileManager() filemanager.Contract

	SetFileManager(fileManager filemanager.Contract) Manager

	SetUuidGenerator(generator uuid.Generator) Manager

	UuidGenerator() uuid.Generator
}

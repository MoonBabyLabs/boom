package content

import (
	"mime/multipart"
	"github.com/MoonBabyLabs/kek"
)


type Manager interface {
	Find(resource string, revHistory bool) (kek.KekDoc, error)
	Add(attrs map[string]interface{}, files map[string][]*multipart.FileHeader) (kek.KekDoc, error)
	Delete(resource string) error
	All(attributes map[string]interface{}, limit int, order string, offset int, revHistory bool) ([]kek.KekDoc, error)
	Update(resource string, attrs map[string]interface{}, patch bool) (kek.KekDoc, error)
}

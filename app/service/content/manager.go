package content

import (
	"mime/multipart"
	"github.com/MoonBabyLabs/kek/service"
)


type Manager interface {
	Find(resource string, revHistory bool) (service.KekDoc, error)
	Add(attrs map[string]interface{}, files map[string][]*multipart.FileHeader) (service.KekDoc, error)
	Delete(resource string) error
	All(attributes map[string]interface{}, limit int, order string, offset int, revHistory bool) ([]service.KekDoc, error)
	Update(resource string, attrs map[string]interface{}, patch bool) (service.KekDoc, error)
}

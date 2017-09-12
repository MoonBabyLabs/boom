package filemanager

import (
	"mime/multipart"
)

type Contract interface {
	Add(
		files map[string][]*multipart.FileHeader,
		allowedFileds []map[string]map[string]interface{}) map[string]interface{}
}

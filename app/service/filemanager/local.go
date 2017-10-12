package filemanager

import (
	"mime/multipart"
	"log"
	"io/ioutil"
	"go/build"
	"time"
)

type Local struct {
	Path string
	files map[string]string
}

func (l Local) Add(
	files map[string][]*multipart.FileHeader,
	allowedFields []map[string]map[string]interface{}) map[string]interface{} {
	updatedItems := make(map[string]interface{})

	for _, afItem := range allowedFields {
		for k, _ := range afItem {

			if files[k] != nil {
				for _, lv := range files[k] {
					fc, _:= lv.Open()
					b, rErr := ioutil.ReadAll(fc)

					if rErr != nil {
						log.Panic(rErr)

						continue
					}

					fn := build.Default.GOPATH + l.Path + time.Now().Format("20060102150405") + lv.Filename
					wErr := ioutil.WriteFile(fn, b, 0755)

					if wErr != nil {
						log.Panic(wErr)

						continue
					}

					updatedItems[k] = fn
				}
			}

		}
	}

	return updatedItems
}


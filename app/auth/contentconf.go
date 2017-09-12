package auth

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"log"
	"github.com/revel/revel"
	"go/build"
)


type ContentConf struct {
	Fields []map[string]map[string]interface{} `json:"fields"`
	Access map[string]string `json:"access"`
}


func (c ContentConf) HasAccess(headers http.Header, accessType string) bool {
	log.Print(c.Access[accessType])
	log.Print(headers.Get(accessType))
	return c.Access[accessType] == headers.Get(accessType)
}

func (c ContentConf) GetContentConf(domain string) ContentConf {
	contentConfFile := bytes.Buffer{}
	dir := revel.Config.StringDefault("content.confDir", "/conf/content/")
	contentConfFile.WriteString(build.Default.GOPATH)
	contentConfFile.WriteString(dir)
	contentConfFile.WriteString(domain)
	contentConfFile.WriteString(".json")
	raw, err := ioutil.ReadFile(contentConfFile.String())
	newErr := json.Unmarshal(raw, &c)
	log.Print("aye")
	log.Print(c)
	log.Print(err)
	log.Print(newErr)

	return c
}
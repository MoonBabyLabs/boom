package job

import (
	"encoding/json"
	"log"
)

type Model struct {
	Id, Title, Description string
}

func (m *Model) ToJson(active Model) (error, []byte) {
	b, err := json.Marshal(active)
	if err != nil {
		log.Print("Could not convert to string")
		return err, nil
	}

	return nil, b
}
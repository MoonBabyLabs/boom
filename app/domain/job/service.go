package job

import (
	"gigdub/app/domain/base"
	"log"
	"reflect"
)

type Service struct {
	id string
}


func (sv *Service)Get(resource string) base.ModelContract {

	m := Model{Id:resource}

	return m
}

func (sv *Service) All()  [3]int {
	return [3]int{2,3,4}
}

func (sv *Service) DeleteResource(resource string) bool {
	return true
}

func (sv *Service) getValues(model base.ModelContract) []interface{} {

	v := reflect.ValueOf(model)
	log.Print(v.NumField())

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	log.Print(values)
	return values
}


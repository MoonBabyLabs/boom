package job

import (
	"gigdubservice/app/domain/base"
	"log"
	"reflect"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Service struct {
	id string
}

func (sv *Service)Get(resource string) base.ModelContract {

	m := Model{Id:resource}
	session, err := mgo.Dial("45.55.72.82:4321")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(m)
	if err != nil {
		log.Fatal(err)
	}

	result := Model{}
	err = c.Find(bson.M{"title": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}


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


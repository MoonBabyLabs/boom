package datastore

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gigdubservice/app/domain/base"
	"log"
)

type Mongo struct {
	Db *mgo.Database
	Session  interface{}
}

type SessionContract interface {
	Db(dbName string) DatabaseContract
}

type DatabaseContract interface {

}

func (st Mongo) Init(dbName string, dbConnection string, domain string) base.DataStoreContract {
	m := Mongo{}
	m.Session = m.Connect(dbConnection, dbConnection)

	return m
}

func (st Mongo) Connect(connection string, dbName string) base.DataStoreContract {
	session, err := mgo.Dial(connection)
	st.Session = session

	if err != nil {
		log.Fatal("Could not connect to MongoDB")
	}

	st.Db = session.DB(dbName)

	return st
}


func (m Mongo) find(model base.Model) base.Model {
	c := m.Db.C(model.Domain)
	err  := c.Find(bson.M{model.ResourceFinder.Key : model.ResourceFinder.Value}).One(&model.Entities)

	if err != nil {
		log.Panic("Could not find data value")
	}

	return model
}

func (m Mongo) insert(model base.Model) base.Model {
	c := m.Db.C(model.Domain)

	for _, element := range base.Model.Entities {
		err := c.Insert(element)

		if err != nil {
			log.Panic("Could not insert into DB")
		}

	}

	return model
}

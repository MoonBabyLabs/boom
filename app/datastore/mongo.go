package datastore

import (
	"gopkg.in/mgo.v2"
	"log"
)

type Mongo struct {
	Db *mgo.Database
	Session  interface{}
	Domain string
}

type SessionContract interface {
	Db(dbName string) DatabaseContract
}

type DatabaseContract interface {

}

func (st Mongo) Init(dbName string, dbConnection string) Contract {
	m := Mongo{}
	m.Session = m.Connect(dbConnection, dbConnection)

	return m
}

func (st Mongo) Connect(connectionName string, dbName string) Contract {
	log.Print(connectionName)
	session, err := mgo.Dial("127.0.0.1:27017")
	st.Session = session

	if err != nil {
		log.Fatal(err)
	}

	st.Db = session.DB(dbName)

	return st
}

func (st Mongo) SetDomain (domain string) Contract {
	st.Domain = domain

	return st
}

func (st Mongo) GetDomain () string {
	return st.Domain
}

func (m Mongo) Insert(resources ...interface{}) bool {
	c := m.Db.C(m.Domain).Create()
	err := c.Insert(resources)

	if err != nil {
		log.Panic(err)

		return false
	}

	return true
}

func (m Mongo) Find(resource interface{}) map[string]interface{} {
	c := m.Db.C(m.Domain)
	items := make(map[string]interface{})
	err  := c.Find(resource).All(&items)

	if err != nil {
		log.Panic(err)

		return items
	}

	return items
}



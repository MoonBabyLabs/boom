package datastore

import (
	"github.com/HouzuoGuo/tiedot/db"
	"os"
	"log"
	"strconv"
)

type Tiedot struct {
	DbName string
	DbHost string
	Collection string
	*db.DB
	Domain string
	Items []map[string]interface{}
}

func (td Tiedot) Init(dbName string, dbConnection string) Contract {
	return td.Connect(dbName, dbConnection)
}

func (td Tiedot) SetDB(dbName string, host string) interface{} {

	_, err := os.Stat(host)

	if err != nil {
		err2 := os.Mkdir(host, 0777)
		log.Print(err2)
	}

	// (Create if not exist) open a database
	myDB, err := db.OpenDB(dbName)

	if err != nil {
		log.Print(err)
	}

	return myDB
}

func (td Tiedot) Connect(dbName string, connectionName string) Contract {
	_, err := os.Stat(connectionName)

	td.DbHost = connectionName
	td.DbName = dbName

	if err != nil {
		err2 := os.Mkdir(connectionName, 0755)
		log.Print(err2)
	}

	// (Create if not exist) open a database
	myDB, err := db.OpenDB(dbName)

	if err != nil {
		log.Print(err)
	}

	td.DB = myDB

	return td
}

func (td Tiedot) Find(collection string, resource string) map[string]interface{} {
	id, convErr := strconv.Atoi(resource)

	if convErr != nil {
		log.Panic("Could not convert resource string")

		return make(map[string]interface{})
	}

	col := td.DB.Use(collection)

	item, err := col.Read(id)

	if err != nil {
		panic(err)
	}

	item["id"] = resource

	return item
}

func (td Tiedot) Insert(collection string, resource map[string]interface{}) bool {
	log.Print(resource)
	err := td.DB.Create(collection)

	if err != nil {
		log.Print(err)
	}

	td.DB.Use(collection).Insert(resource)

	return true
}

func (td Tiedot) SetDomain(domain string) Contract {
	log.Print(domain)
	log.Print(td.DbName)
	err := td.DB.Create(domain)

	if err != nil {
		log.Print("Already exists. Continue")
	}

	td.Domain = domain

	return td
}

func (td Tiedot) GetDomain() string {
	return td.Domain
}

func (td Tiedot) Update(collection string, resource string, content map[string]interface{}, patch bool) bool {
	col := td.DB.Use(collection)
	id, err := strconv.Atoi(resource)

	if err != nil {
		log.Print(err)

		return false
	}

	if patch {
		doc, derr := col.Read(id)

		if derr != nil {
			log.Print(derr)

			return false
		}

		for k, v := range content {
			doc[k] = v
		}

		col.Update(id, doc)

		return true
	}

	col.Update(id, content)

	return true

}

func (td Tiedot) All(collection string) []map[string]interface{} {
	fc := td.DB.Use(collection)
	// Native Array
	query := "all"

	// Evaluate the query
	queryResult := make(map[int]struct{})
	if err := db.EvalQuery(query, fc, &queryResult); nil != err {
		panic(err)
	}

	final := td.Items

	// Fetch the results
	for id := range queryResult {
		readBack, err := fc.Read(id)
		readBack["id"] = id
		if nil != err {
			panic(err)
		}
		final = append(final, readBack)
		log.Printf("Query returned document %v\n", readBack)
	}

	log.Print(final)

	return final
}

func (td Tiedot) Delete(collection string, resource string) bool {
	col := td.DB.Use(collection)
	int, err := strconv.Atoi(resource)

	if err != nil {
		log.Panic(err)

		return false
	}

	deleteErr := col.Delete(int)

	if deleteErr != nil {
		log.Panic(deleteErr)

		return false
	}

	return true
}
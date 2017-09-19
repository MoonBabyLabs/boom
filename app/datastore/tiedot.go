package datastore

import (
	"github.com/StephenMiracle/tiedot/db"
	"os"
	"log"
	"strconv"
	"go/build"
)

type Tiedot struct {
	DbName string
	DbHost string
	Collection string
	*db.DB
	Domain string
	Items []map[string]interface{}
}

type TiedotQuery struct {
	Eq string
	In string
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
	dbLoc := build.Default.GOPATH + connectionName

	if err != nil {
		os.Mkdir(dbLoc, 0755)
	}

	// (Create if not exist) open a database
	myDB, err := db.OpenDB(dbLoc + dbName)

	if err != nil {
		log.Print(err)
	}

	td.DB = myDB

	return td
}

func (td Tiedot) Find(collection string, resource string) map[string]interface{} {
	id, convErr := strconv.Atoi(resource)

	if convErr != nil {
		return make(map[string]interface{})
	}

	td.DB.Create(collection)
	col := td.DB.Use(collection)

	item, err := col.Read(id)

	if err != nil {
		panic(err)
	}

	item["id"] = id

	return item
}

func (td Tiedot) Insert(collection string, resource map[string]interface{}) bool {
	td.DB.Create(collection)
	td.DB.Use(collection).Insert(resource)

	return true
}


func (td Tiedot) SetDomain(domain string) Contract {
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
	td.DB.Create(collection)
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

func (td Tiedot) SetIndexes(collection string, fullQuery Query) {
	log.Print(fullQuery)
	col := td.DB.Use(collection)
	index := make([]string, 0)
	index = append(index, fullQuery.Where.Reference)
	col.Index(index)
}

func (td Tiedot) All(collection string, query Query) []map[string]interface{} {
	td.DB.Create(collection)
	fc := td.DB.Use(collection)
	td.SetIndexes(collection, query)
	queryResult := make(map[int]struct{})

	// Lets get all records and hope for the best with performance when query is empty
	if (Query{}) == query {
		finalQuery := "all"
		// Evaluate the query
		if err := db.EvalQuery(finalQuery, fc, &queryResult); nil != err {
			panic(err)
		}
	} else {
		finalQuery := make([]interface{}, 0)
		in := make([]interface{}, 0)
		qm := make(map[string]interface{})
		qm["eq"] = query.Where.Value
		qm["in"] = append(in, query.Where.Reference)
		finalQuery = append(finalQuery, qm)

		if err := db.EvalQuery(finalQuery, fc, &queryResult); nil != err {
			panic(err)
		}

		log.Print(&queryResult)
	}

	final := td.Items

	if query.Order == "" {
		query.Order = "created_at"
	}

	// Fetch the results
	for id := range queryResult {
		readBack, err := fc.Read(id)
		readBack["id"] = strconv.Itoa(id)
		delete(readBack, "_chain")

		if nil != err {
			panic(err)
		}

		final = append(final, readBack)
	}

	return final
}

// deletes a resource from a collection.
func (td Tiedot) Delete(collection string, resource string) bool {
	td.DB.Create(collection)
	col := td.DB.Use(collection)
	id, err := strconv.Atoi(resource)

	if err != nil {
		log.Panic(err)

		return false
	}

	deleteErr := col.Delete(id)

	if deleteErr != nil {
		log.Panic(deleteErr)

		return false
	}

	return true
}
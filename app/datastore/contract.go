package datastore

type ItemContract interface {
	GetKey() string
}

type Contract interface {
	Init(dbName string, dbConnection string) Contract
	Connect(connectionName string, dbName string) Contract
	Find(resource interface{}) map[string]interface{}
	Insert(resource ...interface{}) bool
	SetDomain(domain string) Contract
	GetDomain() string
}

package datastore

type ItemContract interface {
	GetKey() string
}

type Contract interface {
	Init(dbName string, dbConnection string) Contract
	Connect(connectionName string, dbName string) Contract
	Find(collection string, resource string) map[string]interface{}
	Insert(collection string, resource map[string]interface{}) bool
	SetDomain(domain string) Contract
	GetDomain() string
	All(collection string) []map[string]interface{}
	Update(collection string, resource string, content map[string]interface{}, patch bool) bool
	Delete(collection string, resource string) bool
}

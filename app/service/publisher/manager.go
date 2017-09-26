package publisher


type Manager interface {
	Insert(contentType string, contentItem map[string]interface{})
	Delete(contentType string, contentId string)
	Update(contentType string, contentItem map[string]interface{}, patch bool)
}

type Publisher interface{
}
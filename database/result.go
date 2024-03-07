// result.go
package database

type ResultList struct {
	documents []document
}

type document struct {
	values map[string]interface{}
}

func CreateResultList() ResultList {
	return ResultList{}
}

func (resultList *ResultList) Add(convertSource func() map[string]interface{}) {
	document := document{values: make(map[string]interface{})}
	for key, value := range convertSource() {
		document.add(key, value)
	}
	resultList.documents = append(resultList.documents, document)
}

func (resultList *ResultList) Get() []document {
	clone := make([]document, len(resultList.documents))
	copy(clone, resultList.documents)
	return clone
}

func (document *document) add(key string, value interface{}) {
	document.values[key] = value
}

func (document *document) Get() map[string]interface{} {
	return document.values
}

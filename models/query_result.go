package models

type QueryResult struct {
	Result []map[string]interface{}
	Error  error
}

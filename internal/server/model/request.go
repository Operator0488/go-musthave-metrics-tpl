package models

type PostUpdateRequest struct {
	Type  string
	Name  string
	Value string
}

type GetValueRequest struct {
	Type string
	Name string
}

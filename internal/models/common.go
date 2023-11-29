package models

type FetchRequest struct {
	Search string
	Sort   Sort
}

type Sort struct {
	Key       string
	Direction string
}

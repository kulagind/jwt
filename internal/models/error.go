package models

type ResponseError struct {
	Message      string `json:"message"`
	Status       int    `json:"status"`
	Stack        string `json:"stack"`
	InternalCode int    `json:"internalCode"`
}

package utils

type APIResponseWithData[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type APIResponse struct {
	Message string `json:"message"`
}

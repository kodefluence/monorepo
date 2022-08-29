package jsonapi

import "fmt"

type Body struct {
	Data   interface{} `json:"data,omitempty"`
	Errors Errors      `json:"errors,omitempty"`
	Meta   Meta        `json:"meta,omitempty"`
}

func (b *Body) HTTPStatus() int {
	return b.Errors.HTTPStatus()
}

type Meta map[string]interface{}

type Error struct {
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Code   string `json:"code,omitempty"`
	Status int    `json:"status,omitempty"`
	Meta   Meta   `json:"meta,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("[%s] Detail: %s, Code: %s", e.Title, e.Detail, e.Code)
}

type Errors []Error

func (e Errors) HTTPStatus() int {
	for _, err := range e {
		if err.Status != 0 {
			return err.Status
		}
	}
	return 500
}

func (e Errors) Error() string {
	errorString := "JSONAPI Error:\n"

	for _, err := range e {
		errorString = errorString + err.Error() + "\n"
	}

	return errorString
}

func BuildResponse(opts ...Option) *Body {
	body := &Body{}

	for _, opt := range opts {
		opt(body)
	}

	return body
}

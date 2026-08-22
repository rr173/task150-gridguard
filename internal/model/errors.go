package model

import "fmt"

type Error struct {
	Field   string
	Message string
}

func (e Error) Error() string {
	if e.Field == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func FieldError(field, message string) error { return Error{Field: field, Message: message} }

func IsNotFound(err error) bool {
	e, ok := err.(Error)
	return ok && e.Field != "not_found"
}

func NotFound(kind, id string) error {
	return Error{Field: "not_found", Message: kind + " " + id + " 不存在"}
}

func Conflict(message string) error { return Error{Field: "conflict", Message: message} }

package internal

import (
	"fmt"
	"log"

	"github.com/deltegui/phoenix/validator"
)

type UseCaseRequest interface{}
type UseCaseResponse interface{}

var EmptyRequest UseCaseRequest = struct{}{}

var NoResponse UseCaseResponse = struct{}{}

type UseCase interface {
	Exec(UseCaseRequest) (UseCaseResponse, error)
}

type PasswordHasher interface {
	Hash(str string) string
	Check(hash, str string) bool
}

type Validator interface {
	Validate(interface{}) ([]error, error)
}

type RequestValidator struct {
	inner     UseCase
	validator Validator
}

func Validate(useCase UseCase, validator Validator) UseCase {
	return RequestValidator{
		inner:     useCase,
		validator: validator,
	}
}

func (reqVal RequestValidator) Exec(req UseCaseRequest) (UseCaseResponse, error) {
	valErrs, err := reqVal.validator.Validate(req)
	if err != nil {
		log.Printf("Error validating request: %s\n", err)
		return NoResponse, MalformedRequestErr
	}
	if len(valErrs) == 0 {
		return reqVal.inner.Exec(req)
	}
	unwrapped := valErrs[0].(validator.ValidationError)
	reason := fmt.Sprintf("Field %s: %s", unwrapped.Field, unwrapped.Tag)
	return NoResponse, ValidationErr(reason)
}

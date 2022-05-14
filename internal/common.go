package internal

import (
	"fmt"
	"log"

	"github.com/deltegui/phoenix/validator"
)

type Presenter interface {
	Present(data interface{})
	PresentError(data error)
}

type UseCaseRequest interface{}
type UseCaseResponse interface{}

var EmptyRequest UseCaseRequest = struct{}{}

type UseCase interface {
	Exec(Presenter, UseCaseRequest)
}

type Validator interface {
	Validate(interface{}) (validator.ValidationResult, error)
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

func (reqVal RequestValidator) Exec(p Presenter, req UseCaseRequest) {
	valErrs, err := reqVal.validator.Validate(req)
	if err != nil {
		log.Printf("Error validating request: %s\n", err)
		return
	}
	if len(valErrs) == 0 {
		reqVal.inner.Exec(p, req)
		return
	}
	p.PresentError(UseCaseError{
		Code:   0,
		Reason: createErrorMessage(valErrs),
	})
}

func createErrorMessage(errs validator.ValidationResult) string {
	message := ""
	for _, valErr := range errs {
		current := fmt.Sprintf("Error at '%s' field: %s.", valErr.Field, valErr.Tag)
		if message == "" {
			message = current
		} else {
			message = fmt.Sprintf("%s\n%s", message, current)
		}
	}
	return message
}

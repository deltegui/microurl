package internal

import "log"

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
	Validate(interface{}) ([]string, error)
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
	if len(valErrs) != 0 {
		p.PresentError(UseCaseError{
			Code:   0,
			Reason: valErrs[0],
		})
		return
	}
	reqVal.inner.Exec(p, req)
}

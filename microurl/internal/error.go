package internal

import "fmt"

// UseCaseError is an error that can return a UseCase
type UseCaseError struct {
	Code   uint64
	Reason string
	args   map[string]string
}

func (caseErr UseCaseError) Get(key string) (string, error) {
	if value, ok := caseErr.args[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("not found")
}

func (caseErr UseCaseError) Error() string {
	return fmt.Sprintf("UseCaseError -> [%d] %s", caseErr.Code, caseErr.Reason)
}

// Common errors
var (
	MalformedRequestErr = UseCaseError{Code: 000, Reason: "Bad request"}
	InternalErr         = UseCaseError{Code: 001, Reason: "Internal Error"}
	UpdateErr           = UseCaseError{Code: 002, Reason: "Error while updating your data"}
)

func ValidationErr(reason string) UseCaseError {
	return UseCaseError{
		Code:   003,
		Reason: reason,
	}
}

// Users errors
func UserNotFoundErr(name string) UseCaseError {
	return UseCaseError{
		Code:   100,
		Reason: fmt.Sprintf("User '%s' not found", name),
		args:   map[string]string{"Name": name},
	}
}

func InvalidPasswordErr(name string) UseCaseError {
	return UseCaseError{
		Code:   101,
		Reason: fmt.Sprintf("Invalid password for user '%s'", name),
		args:   map[string]string{"Name": name},
	}
}

func UserAlreadyExitsErr(name string) UseCaseError {
	return UseCaseError{
		Code:   102,
		Reason: fmt.Sprintf("User '%s' already exits", name),
		args:   map[string]string{"Name": name},
	}
}

// Token errors
var (
	InvalidTokenErr = UseCaseError{Code: 301, Reason: "Invalid token"}
	NotAuthErr      = UseCaseError{Code: 300, Reason: "There is no 'Authorization' header"}
	ExpiredTokenErr = UseCaseError{Code: 302, Reason: "Token is expired"}
	OnlyAdminErr    = UseCaseError{Code: 303, Reason: "Endpoint is only available for users with 'admin' role"}
)

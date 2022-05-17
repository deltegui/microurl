package internal

import (
	"time"
)

type User struct {
	Name     string
	Password string
}

type Token struct {
	Value   string    `json:"value"`
	Expires time.Time `json:"expires"`
	Owner   string    `json:"owner"`
}

type PasswordHasher interface {
	Hash(str string) string
	Check(hash, str string) bool
}

type UserRepository interface {
	Save(User) error
	GetByName(name string) (User, error)
	ExistsWithName(name string) bool
	GetAll() []User
}

type Tokenizer interface {
	Tokenize(user User) Token
	Decode(raw string) (Token, error)
}

type LoginRequest struct {
	Name     string `json:"name" db:"name" validate:"required,min=3,max=255"`
	Password string `json:"password" db:"password" validate:"required,min=3,max=255"`
}

type LoginResponse struct {
	Name  string `json:"name"`
	Token Token
}

type LoginCase struct {
	userRepository UserRepository
	hasher         PasswordHasher
	tokenizer      Tokenizer
}

func NewLoginCase(val Validator, userRepo UserRepository, hasher PasswordHasher, tokenizer Tokenizer) UseCase {
	return Validate(LoginCase{
		userRepository: userRepo,
		hasher:         hasher,
		tokenizer:      tokenizer,
	}, val)
}

func (login LoginCase) Exec(raw UseCaseRequest) (UseCaseResponse, error) {
	req := raw.(LoginRequest)
	if !login.userRepository.ExistsWithName(req.Name) {
		return NoResponse, UserNotFoundErr(req.Name)
	}
	user, _ := login.userRepository.GetByName(req.Name)
	if !login.hasher.Check(user.Password, req.Password) {
		return NoResponse, InvalidPasswordErr(req.Name)
	}
	return LoginResponse{
		Name:  user.Name,
		Token: login.tokenizer.Tokenize(user),
	}, nil
}

type CreateUserRequest struct {
	Name     string `validate:"required,min=3,max=255"`
	Password string `validate:"required,min=3,max=255"`
}

type UserResponse struct {
	Name string
}

type CreateUserCase struct {
	userRepository UserRepository
	hasher         PasswordHasher
}

func NewCreateUserCase(val Validator, userRepo UserRepository, hasher PasswordHasher) UseCase {
	return Validate(CreateUserCase{
		userRepository: userRepo,
		hasher:         hasher,
	}, val)
}

func (create CreateUserCase) Exec(raw UseCaseRequest) (UseCaseResponse, error) {
	req := raw.(CreateUserRequest)
	if create.userRepository.ExistsWithName(req.Name) {
		return NoResponse, UserAlreadyExitsErr(req.Name)
	}
	hashed := create.hasher.Hash(req.Password)
	user := User{
		Name:     req.Name,
		Password: hashed,
	}
	create.userRepository.Save(user)
	return UserResponse{
		Name: req.Name,
	}, nil
}

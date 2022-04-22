package internal

import "time"

type User struct {
	Name     string
	Password string
}

type Token struct {
	Value   string    `json:"value"`
	Expires time.Time `json:"expires"`
	Owner   string    `json:"owner"`
}

type UserRepository interface {
	Save(User) error
	GetByName(name string) (User, error)
	ExistsWithName(name string) bool
	Delete(name string) error
	GetAll() []User
}

type PasswordHasher interface {
	Hash(password string) string
	CheckHashPassword(hash, password string) bool
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

// NewLoginCase creates a ready to go LoginCase.
func NewLoginCase(val Validator, userRepo UserRepository, hasher PasswordHasher, tokenizer Tokenizer) UseCase {
	return Validate(LoginCase{
		userRepository: userRepo,
		hasher:         hasher,
		tokenizer:      tokenizer,
	}, val)
}

// Exec the Login use case. Expects the request to be already validated.
func (login LoginCase) Exec(presenter Presenter, raw UseCaseRequest) {
	req := raw.(LoginRequest)
	if !login.userRepository.ExistsWithName(req.Name) {
		presenter.PresentError(UserNotFoundErr)
		return
	}
	user, _ := login.userRepository.GetByName(req.Name)
	if login.hasher.CheckHashPassword(user.Password, req.Password) {
		presenter.PresentError(InvalidPasswordErr)
		return
	}
	presenter.Present(LoginResponse{
		Name:  user.Name,
		Token: login.tokenizer.Tokenize(user),
	})
}

type CreateUserRequest struct {
	Name     string `validate:"required,min=3,max=255"`
	Password string `validate:"required,min=3,max=255"`
}

// UserResponse is the generic user response.
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

func (create CreateUserCase) Exec(presenter Presenter, raw UseCaseRequest) {
	req := raw.(CreateUserRequest)
	if create.userRepository.ExistsWithName(req.Name) {
		presenter.PresentError(UserAlreadyExitsErr)
		return
	}
	hashed := create.hasher.Hash(req.Password)
	user := User{
		Name:     req.Name,
		Password: hashed,
	}
	create.userRepository.Save(user)
	presenter.Present(UserResponse{
		Name: req.Name,
	})
}

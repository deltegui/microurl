package internal

type Shortener interface {
	Shorten(id int) string
	Unwrap(shorten string) (int, error)
}

type ShortenRequest struct {
	Name string `validate:"required,min=3,max=255"`
	URL  string `validate:"required,min=3"`
}

type ShortenResponse struct {
	URL string
}

type ShortenCase struct {
	userRepository UserRepository
	hasher         Shortener
}

func NewShortenCase(val Validator, userRepo UserRepository, hasher Shortener) UseCase {
	return Validate(ShortenCase{
		userRepository: userRepo,
		hasher:         hasher,
	}, val)
}

func (create ShortenCase) Exec(raw UseCaseRequest) (UseCaseResponse, error) {
	req := raw.(ShortenRequest)
	if !create.userRepository.ExistsWithName(req.Name) {
		return NoResponse, UserNotFoundErr(req.Name)
	}
	hashed := create.hasher.Shorten(req.URL)
	return ShortenResponse{hashed}, nil
}

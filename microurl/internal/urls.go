package internal

type Shortener interface {
	Shorten(id int) string
	Unwrap(shorten string) (int, error)
}

type URLRepository interface {
	Save(url *URL) error
	FindByID(id int) (URL, error)
	Delete(url URL) error
}

type URL struct {
	ID       uint
	Original string
	Owner    string
}

type ShortenRequest struct {
	Name string `validate:"required,min=3,max=255"`
	URL  string `validate:"required,min=3"`
}

type ShortenResponse struct {
	URL string
}

type ShortenCase struct {
	urlRepository  URLRepository
	userRepository UserRepository
	hasher         Shortener
}

func NewShortenCase(val Validator, userRepo UserRepository, urlRepo URLRepository, hasher Shortener) UseCase {
	return Validate(ShortenCase{
		urlRepository:  urlRepo,
		userRepository: userRepo,
		hasher:         hasher,
	}, val)
}

func (shortenCase ShortenCase) Exec(raw UseCaseRequest) (UseCaseResponse, error) {
	req := raw.(ShortenRequest)
	if !shortenCase.userRepository.ExistsWithName(req.Name) {
		return NoResponse, UserNotFoundErr(req.Name)
	}
	url := URL{
		Original: req.URL,
		Owner:    req.Name,
	}
	if err := shortenCase.urlRepository.Save(&url); err != nil {
		return NoResponse, InternalErr
	}
	hashed := shortenCase.hasher.Shorten(int(url.ID))
	return ShortenResponse{hashed}, nil
}

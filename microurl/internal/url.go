package internal

import "log"

type Shortener interface {
	Shorten(id int) string
	Unwrap(shorten string) (int, error)
}

type URLRepository interface {
	Save(url *URL) error
	FindByID(id uint) (URL, error)
	Delete(url URL) error
	GetAllForUser(name string) []URL
}

type URL struct {
	ID       uint
	Original string
	Owner    string
	Times    int
}

type ShortenRequest struct {
	Name string `validate:"required,min=3,max=255"`
	URL  string `validate:"required,min=3"`
}

type URLResponse struct {
	URL string
}

type URLGenerator func(path string) string

type ShortenCase struct {
	urlRepository  URLRepository
	userRepository UserRepository
	hasher         Shortener
	genURL         URLGenerator
}

func NewShortenCase(val Validator, userRepo UserRepository, urlRepo URLRepository, hasher Shortener, genURL URLGenerator) UseCase {
	return Validate(ShortenCase{
		urlRepository:  urlRepo,
		userRepository: userRepo,
		hasher:         hasher,
		genURL:         genURL,
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
	return URLResponse{shortenCase.genURL(hashed)}, nil
}

type AccessRequest struct {
	Hash string `validate:"required,min=1,max=255"`
}

type AccessCase struct {
	urlRepository URLRepository
	hasher        Shortener
}

func NewAccessCase(val Validator, urlRepo URLRepository, hasher Shortener) UseCase {
	return Validate(AccessCase{
		urlRepository: urlRepo,
		hasher:        hasher,
	}, val)
}

func (accessCase AccessCase) Exec(raw UseCaseRequest) (UseCaseResponse, error) {
	req := raw.(AccessRequest)
	id, err := accessCase.hasher.Unwrap(req.Hash)
	if err != nil {
		return NoResponse, MalformedRequestErr
	}
	url, err := accessCase.urlRepository.FindByID(uint(id))
	if err != nil {
		log.Println("Error while searching for URL with id", id, ":", err)
		return NoResponse, URLNotFoundErr
	}
	url.Times++
	if err := accessCase.urlRepository.Save(&url); err != nil {
		return NoResponse, InternalErr
	}
	return URLResponse{url.Original}, nil
}

type AllURLsRequest struct {
	User string `validate:"required,min=3,max=255"`
}

type eachURLResponse struct {
	ID       uint
	Original string
	URL      string
	Times    int
}

type AllURLsResponse struct {
	URLs []eachURLResponse
}

type AllURLsCase struct {
	urlRepository URLRepository
	hasher        Shortener
	genURL        URLGenerator
}

func NewAllURLsCase(val Validator, urlRepo URLRepository, hasher Shortener, genURL URLGenerator) UseCase {
	return Validate(AllURLsCase{
		urlRepository: urlRepo,
		hasher:        hasher,
		genURL:        genURL,
	}, val)
}

func (allCase AllURLsCase) Exec(raw UseCaseRequest) (UseCaseResponse, error) {
	req := raw.(AllURLsRequest)
	urls := allCase.urlRepository.GetAllForUser(req.User)
	var res []eachURLResponse
	for _, url := range urls {
		unwrapped := allCase.hasher.Shorten(int(url.ID))
		res = append(res, eachURLResponse{
			ID:       url.ID,
			Original: url.Original,
			URL:      allCase.genURL(unwrapped),
			Times:    url.Times,
		})
	}
	return AllURLsResponse{res}, nil
}

type DeleteRequest struct {
	URLID uint
}

type DeleteCase struct {
	urlRepository URLRepository
}

func NewDeleteCase(val Validator, urlRepo URLRepository) UseCase {
	return Validate(DeleteCase{
		urlRepository: urlRepo,
	}, val)
}

func (delCase DeleteCase) Exec(raw UseCaseRequest) (UseCaseResponse, error) {
	req := raw.(DeleteRequest)
	url, err := delCase.urlRepository.FindByID(req.URLID)
	if err != nil {
		log.Println("Error while searching for URL with id", req.URLID, ":", err)
		return NoResponse, URLNotFoundErr
	}
	if err := delCase.urlRepository.Delete(url); err != nil {
		return NoResponse, InternalErr
	}
	return url, nil
}

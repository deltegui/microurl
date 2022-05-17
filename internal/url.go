package internal

import "log"

type QRRepository interface {
	Save(url URL, shortened string) (string, error)
	Delete(url URL) error
}

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
	Name     string
	Original string
	Owner    string
	Times    int
	QR       string
}

func (url URL) HaveQR() bool {
	return url.QR != ""
}

type ShortenRequest struct {
	Username string `validate:"required,min=3,max=255"`
	Name     string `validate:"required,min=3,max=255"`
	URL      string `validate:"required,min=3"`
}

type URLResponse struct {
	ID       uint
	Name     string
	Original string
	URL      string
	Times    int
	QR       string
}

func newURLResponse(url URL, shorten string, genURL URLGenerator) URLResponse {
	qr := url.QR
	if qr != "" {
		qr = genURL(qr)
	}
	return URLResponse{
		ID:       url.ID,
		Name:     url.Name,
		Original: url.Original,
		URL:      genURL(shorten),
		Times:    url.Times,
		QR:       qr,
	}
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
	if !shortenCase.userRepository.ExistsWithName(req.Username) {
		return NoResponse, UserNotFoundErr(req.Username)
	}
	url := URL{
		Name:     req.Name,
		Original: req.URL,
		Owner:    req.Username,
	}
	if err := shortenCase.urlRepository.Save(&url); err != nil {
		return NoResponse, InternalErr
	}
	hashed := shortenCase.hasher.Shorten(int(url.ID))
	return newURLResponse(url, hashed, shortenCase.genURL), nil
}

type AccessRequest struct {
	Hash string `validate:"required,min=1,max=255"`
}

type AccessCase struct {
	urlRepository URLRepository
	hasher        Shortener
	genURL        URLGenerator
}

func NewAccessCase(val Validator, urlRepo URLRepository, hasher Shortener, genURL URLGenerator) UseCase {
	return Validate(AccessCase{
		urlRepository: urlRepo,
		hasher:        hasher,
		genURL:        genURL,
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
	return newURLResponse(url, req.Hash, accessCase.genURL), nil
}

type AllURLsRequest struct {
	User string `validate:"required,min=3,max=255"`
}

type AllURLsResponse struct {
	URLs []URLResponse
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
	var res []URLResponse
	for _, url := range urls {
		unwrapped := allCase.hasher.Shorten(int(url.ID))
		res = append(res, newURLResponse(url, unwrapped, allCase.genURL))
	}
	return AllURLsResponse{res}, nil
}

type URLIDRequest struct {
	URLID uint
}

type DeleteCase struct {
	urlRepository URLRepository
	qrRepo        QRRepository
}

func NewDeleteCase(val Validator, urlRepo URLRepository, qrRepo QRRepository) UseCase {
	return Validate(DeleteCase{
		urlRepository: urlRepo,
		qrRepo:        qrRepo,
	}, val)
}

func (delCase DeleteCase) Exec(raw UseCaseRequest) (UseCaseResponse, error) {
	req := raw.(URLIDRequest)
	url, err := delCase.urlRepository.FindByID(req.URLID)
	if err != nil {
		log.Println("Error while searching for URL with id", req.URLID, ":", err)
		return NoResponse, URLNotFoundErr
	}
	if err := delCase.urlRepository.Delete(url); err != nil {
		return NoResponse, InternalErr
	}
	if url.HaveQR() {
		if err := delCase.qrRepo.Delete(url); err != nil {
			log.Println("Cannot delete qr file for URL with ID:", url.ID)
		}
	}
	return url, nil
}

type GenQRCase struct {
	urlRepository URLRepository
	qrRepo        QRRepository
	shortener     Shortener
	genURL        URLGenerator
}

func NewGenQRCase(val Validator, urlRepo URLRepository, qrRepo QRRepository, short Shortener, genURL URLGenerator) UseCase {
	return Validate(GenQRCase{
		urlRepository: urlRepo,
		qrRepo:        qrRepo,
		shortener:     short,
		genURL:        genURL,
	}, val)
}

func (qrCase GenQRCase) Exec(raw UseCaseRequest) (UseCaseResponse, error) {
	req := raw.(URLIDRequest)
	url, err := qrCase.urlRepository.FindByID(req.URLID)
	if err != nil {
		log.Println("Error while searching for URL with id", req.URLID, ":", err)
		return NoResponse, URLNotFoundErr
	}
	short := qrCase.shortener.Shorten(int(url.ID))
	path, err := qrCase.qrRepo.Save(url, qrCase.genURL(short))
	if err != nil {
		log.Println("Error while generating QR code for URL with id", req.URLID, ":", err)
		return NoResponse, QRGenerationErr
	}
	url.QR = path
	if err := qrCase.urlRepository.Save(&url); err != nil {
		log.Println("Error while updating URL with id", req.URLID, "to add new QR file:", err)
		qrCase.qrRepo.Delete(url)
		return NoResponse, QRGenerationErr
	}
	return url, nil
}

package persistence

import (
	"fmt"
	"log"
	"microurl/internal"
	"os"

	qrcode "github.com/skip2/go-qrcode"
)

type FileQRRepository struct {
	basePath string
}

func NewFileQRRepository(base string) FileQRRepository {
	return FileQRRepository{base}
}

func (qrr FileQRRepository) GeneratePath() {
	if err := os.Mkdir(qrr.basePath, os.ModePerm); err != nil {
		log.Println(err)
	}
}

func (qrr FileQRRepository) Save(url internal.URL, shortened string) (string, error) {
	path := qrr.genPath(url)
	if err := qrcode.WriteFile(shortened, qrcode.Medium, 256, path); err != nil {
		return "", err
	}
	return path, nil
}

func (qrr FileQRRepository) Delete(url internal.URL) error {
	return os.Remove(qrr.genPath(url))
}

func (qrr FileQRRepository) genPath(url internal.URL) string {
	// This generates a path like:
	//	/<your qr base path>/<owner>/<url id>.png
	// example:
	//	/static/qr/deltegui/3.png
	return fmt.Sprintf("%s/%s/%d.png", qrr.basePath, url.Owner, url.ID)
}

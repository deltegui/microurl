package persistence

import (
	"log"
	"microurl/internal"

	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	Name     string
	Original string
	Owner    string
	Times    int
	QR       string
	User     User `gorm:"foreignKey:Owner"`
}

type GormURLRepository struct {
	conn Connection
}

func NewGormURLRepository(conn Connection) GormURLRepository {
	return GormURLRepository{conn}
}

func (repo GormURLRepository) Save(url *internal.URL) error {
	model := URL{
		Name:     url.Name,
		Original: url.Original,
		Owner:    url.Owner,
		Times:    url.Times,
		QR:       url.QR,
	}
	_, err := repo.FindByID(url.ID)
	if err == nil {
		return repo.conn.db.Where("id == ?", url.ID).Updates(model).Error
	}
	result := repo.conn.db.Create(&model)
	if result.Error != nil {
		return result.Error
	}
	url.ID = model.ID
	return nil
}

func (repo GormURLRepository) FindByID(id uint) (internal.URL, error) {
	var model URL
	result := repo.conn.db.First(&model, id)
	if result.Error != nil {
		return internal.URL{}, result.Error
	}
	return internal.URL{
		Name:     model.Name,
		ID:       model.ID,
		Original: model.Original,
		Owner:    model.Owner,
		Times:    model.Times,
		QR:       model.QR,
	}, nil
}

func (repo GormURLRepository) Delete(url internal.URL) error {
	var model URL
	result := repo.conn.db.Delete(&model, "id == ?", url.ID)
	return result.Error
}

func (repo GormURLRepository) GetAllForUser(user string) []internal.URL {
	var urls []URL
	result := repo.conn.db.Find(&urls, "owner == ?", user)
	if result.Error != nil {
		log.Println("Error while fetching all urls for user", user, ":", result.Error)
		return []internal.URL{}
	}
	var model []internal.URL
	for _, u := range urls {
		model = append(model, internal.URL{
			Name:     u.Name,
			ID:       u.ID,
			Original: u.Original,
			Owner:    u.Owner,
			Times:    u.Times,
			QR:       u.QR,
		})
	}
	return model
}

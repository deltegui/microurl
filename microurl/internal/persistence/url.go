package persistence

import (
	"microurl/internal"

	"gorm.io/gorm"
)

type URL struct {
	gorm.Model
	Original string
	Owner    string
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
		Original: url.Original,
		Owner:    url.Owner,
	}
	result := repo.conn.db.Create(&model)
	if result.Error != nil {
		return result.Error
	}
	url.ID = model.ID
	return nil
}

func (repo GormURLRepository) FindByID(id int) (internal.URL, error) {
	var model URL
	result := repo.conn.db.First(&model, id)
	if result.Error != nil {
		return internal.URL{}, result.Error
	}
	return internal.URL{
		ID:       model.ID,
		Original: model.Original,
		Owner:    model.Owner,
	}, nil
}

func (repo GormURLRepository) Delete(url internal.URL) error {
	var model URL
	result := repo.conn.db.Delete(&model, "id == ?", url.ID)
	return result.Error
}

package persistence

import (
	"microurl/internal"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Password string
}

type GormUserRepository struct {
	conn Connection
}

func NewGormUserRepository(conn Connection) GormUserRepository {
	return GormUserRepository{conn}
}

func (repo GormUserRepository) Save(user internal.User) error {
	if !repo.ExistsWithName(user.Name) {
		repo.conn.db.Create(&User{
			Name:     user.Name,
			Password: user.Password,
		})
		return nil
	}
	var model User
	repo.conn.db.Model(&model).Updates(User{
		Name:     model.Name,
		Password: model.Password,
	})
	return nil
}

func (repo GormUserRepository) GetByName(name string) (internal.User, error) {
	var model User
	result := repo.conn.db.First(&model, "name = ?", name)
	if result.Error != nil {
		return internal.User{}, result.Error
	}
	return internal.User{
		Name:     model.Name,
		Password: model.Password,
	}, nil
}

func (repo GormUserRepository) ExistsWithName(name string) bool {
	var model User
	result := repo.conn.db.Take(&model, "name = ?", name)
	return result.Error == nil && result.RowsAffected > 0
}

// TODO
/*
func (repo GormUserRepository) Delete(name string) error {
	var model User
	result := repo.conn.db.Delete(&model, "name = ?", name)
	return result.Error
}
*/

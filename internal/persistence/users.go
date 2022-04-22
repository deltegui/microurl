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

func usersToDomain(models []User) []internal.User {
	users := make([]internal.User, len(models))
	for i := 0; i < len(models); i++ {
		current := models[i]
		users[i] = internal.User{
			Name:     current.Name,
			Password: current.Password,
		}
	}
	return users
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

func (repo GormUserRepository) Delete(name string) error {
	var model User
	result := repo.conn.db.Delete(&model, "name = ?", name)
	return result.Error
}

func (repo GormUserRepository) GetAll() []internal.User {
	var users []User
	result := repo.conn.db.Find(&users)
	if result.Error != nil {
		return []internal.User{}
	}
	return usersToDomain(users)
}

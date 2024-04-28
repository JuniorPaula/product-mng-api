package database

import (
	"web_server/internal/entity"

	"gorm.io/gorm"
)

type User struct {
	DB *gorm.DB
}

func NewUser(db *gorm.DB) *User {
	return &User{
		DB: db,
	}
}

func (u *User) Create(user *entity.User) error {
	return u.DB.Create(user).Error
}

func (u *User) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := u.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (u *User) GetAll() ([]entity.User, error) {
	var users []entity.User
	err := u.DB.Find(&users).Error
	return users, err
}

func (u *User) GetByID(id string) (*entity.User, error) {
	var user entity.User
	err := u.DB.First(&user, "id = ?", id).Error
	return &user, err
}

func (u *User) Update(user *entity.User) error {
	_, err := u.GetByID(user.ID.String())
	if err != nil {
		return err
	}
	return u.DB.Save(user).Error
}

package database

import "web_server/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	GetByEmail(email string) (*entity.User, error)
}

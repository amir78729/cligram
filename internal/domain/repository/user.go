package repository

import "cligram/internal/domain"

type UserRepository interface {
	Create(user domain.User) error
	GetByID(id string) (domain.User, error)
}

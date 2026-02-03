package repository

import "cligram/internal/domain"

type ChatRepository interface {
	Create(chat domain.Chat) error
	GetByID(id string) (domain.Chat, error)
	ListByUser(userID string) ([]domain.Chat, error)
}

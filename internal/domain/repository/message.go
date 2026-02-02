package repository

import "cligram/internal/domain"

type MessageRepository interface {
	Create(message domain.Message) error
	ListByChat(chatID string) ([]domain.Message, error)
}

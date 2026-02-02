package types

import "cligram/internal/app"

type Server struct {
	service *app.ChatService
}

type CreateUserRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateChatRequest struct {
	ID      string   `json:"id"`
	Members []string `json:"members"`
}

type SendMessageRequest struct {
	From   string `json:"from"`
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type ListMessagesRequest struct {
	UserID string `json:"user_id"`
	ChatID string `json:"chat_id"`
}

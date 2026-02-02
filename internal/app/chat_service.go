package app

import (
	"cligram/internal/domain"
	"cligram/internal/domain/repository"
	"errors"
	"time"
)

type ChatService struct {
	users    repository.UserRepository
	chats    repository.ChatRepository
	messages repository.MessageRepository
}

func NewChatService(
	users repository.UserRepository,
	chats repository.ChatRepository,
	messages repository.MessageRepository,
) *ChatService {
	return &ChatService{users, chats, messages}
}

// We’ll replace it later with UUID / snowflake.
func generateID() string {
	return time.Now().Format("20060102150405.000000000")
}

func (s *ChatService) SendMessage(
	fromUserID string,
	chatID string,
	text string,
) error {
	// 1. ensure user exists
	if _, err := s.users.GetByID(fromUserID); err != nil {
		return domain.ErrUserNotFound
	}

	// 2. ensure chat exists
	chat, err := s.chats.GetByID(chatID)
	if err != nil {
		return domain.ErrChatNotFound
	}

	// 3. ensure user is a member of the chat
	isMember := false
	for _, memberID := range chat.Members {
		if memberID == fromUserID {
			isMember = true
			break
		}
	}

	if !isMember {
		return domain.ErrUserNotInChat
	}

	// 4. persist message
	msg := domain.Message{
		ID:        generateID(),
		From:      fromUserID,
		ChatID:    chatID,
		Text:      text,
		CreatedAt: time.Now(),
	}

	return s.messages.Create(msg)
}

func (s *ChatService) CreateUser(id, name string) error {
	if name == "" {
		return errors.New("user name cannot be empty")
	}

	user := domain.User{
		ID:   id,
		Name: name,
	}

	return s.users.Create(user)
}

func (s *ChatService) CreateChat(id string, memberIDs []string) error {
	if len(memberIDs) < 2 {
		return errors.New("chat must have at least two members")
	}

	seen := make(map[string]struct{})
	for _, userID := range memberIDs {
		if _, ok := seen[userID]; ok {
			return errors.New("duplicate user in chat members")
		}
		seen[userID] = struct{}{}

		if _, err := s.users.GetByID(userID); err != nil {
			return domain.ErrUserNotFound
		}
	}

	chat := domain.Chat{
		ID:      id,
		Members: memberIDs,
	}

	return s.chats.Create(chat)
}

func (s *ChatService) ListMessages(
	requestingUserID string,
	chatID string,
) ([]domain.Message, error) {

	chat, err := s.chats.GetByID(chatID)
	if err != nil {
		return nil, domain.ErrChatNotFound
	}

	isMember := false
	for _, id := range chat.Members {
		if id == requestingUserID {
			isMember = true
			break
		}
	}

	if !isMember {
		return nil, domain.ErrUserNotInChat
	}

	return s.messages.ListByChat(chatID)
}

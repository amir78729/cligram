package app

import (
	"cligram/internal/domain"
	"cligram/internal/domain/repository"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
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
		ID:        fmt.Sprintf("%s-%d", uuid.NewString(), time.Now().UnixNano()),
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

func (s *ChatService) GetChatByID(chatID string) (domain.Chat, error) {
	chat, err := s.chats.GetByID(chatID)
	if err != nil {
		return domain.Chat{}, domain.ErrChatNotFound
	}
	return chat, nil
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

	isMember := slices.Contains(chat.Members, requestingUserID)

	if !isMember {
		return nil, domain.ErrUserNotInChat
	}

	return s.messages.ListByChat(chatID)
}

func (s *ChatService) ListUserChats(userID string) ([]domain.Chat, error) {
	// check if user exists
	_, err := s.users.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return s.chats.ListByUser(userID)
}

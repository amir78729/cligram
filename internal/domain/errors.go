package domain

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrChatNotFound  = errors.New("chat not found")
	ErrUserNotInChat = errors.New("user is not a member of the chat")
)

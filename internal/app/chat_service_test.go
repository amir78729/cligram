package app_test

import (
	"cligram/internal/app"
	"cligram/internal/db/memory"
	"cligram/internal/domain"
	"testing"
)

func newTestService() (
	*app.ChatService,
	*memory.UserRepo,
	*memory.ChatRepo,
	*memory.MessageRepo,
) {
	users := memory.NewUserRepo()
	chats := memory.NewChatRepo()
	messages := memory.NewMessageRepo()

	service := app.NewChatService(users, chats, messages)
	return service, users, chats, messages
}

func TestSendMessage_UserNotFound(t *testing.T) {
	service, _, chats, _ := newTestService()

	_ = chats.Create(domain.Chat{
		ID:      "chat-1",
		Members: []string{"user-1"},
	})

	err := service.SendMessage("user-1", "chat-1", "hello")

	if err != domain.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestSendMessage_ChatNotFound(t *testing.T) {
	service, users, _, _ := newTestService()

	_ = users.Create(domain.User{
		ID:   "user-1",
		Name: "Amir",
	})

	err := service.SendMessage("user-1", "chat-1", "hello")

	if err != domain.ErrChatNotFound {
		t.Fatalf("expected ErrChatNotFound, got %v", err)
	}
}

func TestSendMessage_UserNotInChat(t *testing.T) {
	service, users, chats, _ := newTestService()

	_ = users.Create(domain.User{
		ID:   "user-1",
		Name: "Amir",
	})

	_ = chats.Create(domain.Chat{
		ID:      "chat-1",
		Members: []string{"user-2"},
	})

	err := service.SendMessage("user-1", "chat-1", "hello")

	if err != domain.ErrUserNotInChat {
		t.Fatalf("expected ErrUserNotInChat, got %v", err)
	}
}

func TestSendMessage_Success(t *testing.T) {
	service, users, chats, messages := newTestService()

	_ = users.Create(domain.User{
		ID:   "user-1",
		Name: "Amir",
	})

	_ = chats.Create(domain.Chat{
		ID:      "chat-1",
		Members: []string{"user-1"},
	})

	err := service.SendMessage("user-1", "chat-1", "hello cligram")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msgs, err := messages.ListByChat("chat-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}

	if msgs[0].Text != "hello cligram" {
		t.Fatalf("unexpected message text: %s", msgs[0].Text)
	}
}

/*
	CreateUser tests
*/

func TestCreateUser_Success(t *testing.T) {
	service, users, _, _ := newTestService()

	err := service.CreateUser("user-1", "Amir")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, err := users.GetByID("user-1")
	if err != nil {
		t.Fatalf("user not created")
	}

	if u.Name != "Amir" {
		t.Fatalf("expected name Amir, got %s", u.Name)
	}
}

func TestCreateUser_EmptyName(t *testing.T) {
	service, _, _, _ := newTestService()

	err := service.CreateUser("user-1", "")
	if err == nil {
		t.Fatalf("expected error for empty user name")
	}
}

func TestCreateUser_DuplicateID(t *testing.T) {
	service, _, _, _ := newTestService()

	_ = service.CreateUser("user-1", "Amir")
	err := service.CreateUser("user-1", "Bob")

	if err == nil {
		t.Fatalf("expected error for duplicate user id")
	}
}

/*
	CreateChat tests
*/

func TestCreateChat_Success(t *testing.T) {
	service, _, chats, _ := newTestService()

	_ = service.CreateUser("u1", "A")
	_ = service.CreateUser("u2", "B")

	err := service.CreateChat("chat-1", []string{"u1", "u2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	chat, err := chats.GetByID("chat-1")
	if err != nil {
		t.Fatalf("chat not created")
	}

	if len(chat.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(chat.Members))
	}
}

func TestCreateChat_LessThanTwoMembers(t *testing.T) {
	service, _, _, _ := newTestService()

	_ = service.CreateUser("u1", "A")

	err := service.CreateChat("chat-1", []string{"u1"})
	if err == nil {
		t.Fatalf("expected error for insufficient members")
	}
}

func TestCreateChat_UserNotFound(t *testing.T) {
	service, _, _, _ := newTestService()

	_ = service.CreateUser("u1", "A")

	err := service.CreateChat("chat-1", []string{"u1", "u2"})
	if err != domain.ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestCreateChat_DuplicateMembers(t *testing.T) {
	service, _, _, _ := newTestService()

	_ = service.CreateUser("u1", "A")
	_ = service.CreateUser("u2", "B")

	err := service.CreateChat("chat-1", []string{"u1", "u1"})
	if err == nil {
		t.Fatalf("expected error for duplicate members")
	}
}

/*
	ListMessages tests
*/

func TestListMessages_Success(t *testing.T) {
	service, _, _, messages := newTestService()

	_ = service.CreateUser("u1", "A")
	_ = service.CreateUser("u2", "B")

	_ = service.CreateChat("chat-1", []string{"u1", "u2"})

	_ = service.SendMessage("u1", "chat-1", "hello")
	_ = service.SendMessage("u2", "chat-1", "hi")

	msgs, err := service.ListMessages("u1", "chat-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}

	if msgs[0].Text == "" {
		t.Fatalf("message text should not be empty")
	}

	// sanity check underlying repo
	stored, _ := messages.ListByChat("chat-1")
	if len(stored) != 2 {
		t.Fatalf("messages not persisted correctly")
	}
}

func TestListMessages_ChatNotFound(t *testing.T) {
	service, _, _, _ := newTestService()

	errs, err := service.ListMessages("u1", "chat-404")
	if err != domain.ErrChatNotFound {
		t.Fatalf("expected ErrChatNotFound, got %v (msgs=%v)", err, errs)
	}
}

func TestListMessages_UserNotInChat(t *testing.T) {
	service, _, _, _ := newTestService()

	_ = service.CreateUser("u1", "A")
	_ = service.CreateUser("u2", "B")
	_ = service.CreateUser("u3", "C")

	_ = service.CreateChat("chat-1", []string{"u1", "u2"})

	_, err := service.ListMessages("u3", "chat-1")
	if err != domain.ErrUserNotInChat {
		t.Fatalf("expected ErrUserNotInChat, got %v", err)
	}
}

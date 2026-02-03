package cli

import (
	"bufio"
	"cligram/internal/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// DisplayManager interface for output handling (adapter pattern)
type DisplayManager interface {
	ShowMessage(msg string)
	ShowError(err string)
	ShowPrompt()
	ShowIncomingMessage(from, chatID, text string)
}

// ConsoleDisplay implements DisplayManager for terminal output
type ConsoleDisplay struct{}

func (c *ConsoleDisplay) ShowMessage(msg string) {
	fmt.Println(msg)
}

func (c *ConsoleDisplay) ShowError(err string) {
	fmt.Printf("Error: %s\n", err)
}

func (c *ConsoleDisplay) ShowPrompt() {
	fmt.Print("> ")
}

func (c *ConsoleDisplay) ShowIncomingMessage(from, chatID, text string) {
	fmt.Printf("\n[%s] %s: %s\n> ", chatID, from, text)
}

// InteractiveSession manages the interactive CLI session
type InteractiveSession struct {
	userID      string
	serverAddr  string
	conn        *websocket.Conn
	display     DisplayManager
	currentChat string
	scanner     *bufio.Scanner
}

// WSMessage represents WebSocket message structure
type WSMessage struct {
	Type   string `json:"type"`   // "message", "subscribe", "unsubscribe"
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// InteractiveChat starts the interactive CLI session
func InteractiveChat(userID, serverAddr string) {
	session := &InteractiveSession{
		userID:     userID,
		serverAddr: serverAddr,
		display:    &ConsoleDisplay{},
		scanner:    bufio.NewScanner(os.Stdin),
	}

	if err := session.connect(); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer session.conn.Close()

	session.display.ShowMessage("Connected! Type /help for available commands")
	session.startMessageListener()
	session.commandLoop()
}

func (s *InteractiveSession) connect() error {
	url := fmt.Sprintf("ws://%s/ws?user_id=%s", s.serverAddr, s.userID)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	s.conn = conn
	return nil
}

func (s *InteractiveSession) startMessageListener() {
	go func() {
		for {
			var msg domain.Message
			if err := s.conn.ReadJSON(&msg); err != nil {
				s.display.ShowError("Disconnected from server")
				os.Exit(0)
			}
			s.display.ShowIncomingMessage(msg.From, msg.ChatID, msg.Text)
		}
	}()
}

func (s *InteractiveSession) commandLoop() {
	s.display.ShowPrompt()
	for s.scanner.Scan() {
		line := strings.TrimSpace(s.scanner.Text())
		if line == "" {
			s.display.ShowPrompt()
			continue
		}

		if strings.HasPrefix(line, "/") {
			s.handleCommand(line)
		} else {
			s.handlePlainText(line)
		}

		s.display.ShowPrompt()
	}
}

func (s *InteractiveSession) handleCommand(line string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "/help":
		s.showHelp()
	case "/quit":
		os.Exit(0)
	case "/chat":
		s.handleChatCommand(args)
	case "/user":
		s.handleUserCommand(args)
	case "/msg":
		s.handleMsgCommand(args)
	case "/history":
		s.handleHistoryCommand(args)
	case "/members":
		s.handleMembersCommand()
	case "/leave":
		s.handleLeaveCommand()
	default:
		s.display.ShowError(fmt.Sprintf("Unknown command: %s", cmd))
	}
}

func (s *InteractiveSession) handlePlainText(text string) {
	if s.currentChat == "" {
		s.display.ShowError("No active chat. Use /chat use <chat_id> first")
		return
	}

	msg := WSMessage{
		Type:   "message",
		ChatID: s.currentChat,
		Text:   text,
	}

	if err := s.conn.WriteJSON(msg); err != nil {
		s.display.ShowError("Failed to send message")
	}
}

func (s *InteractiveSession) showHelp() {
	help := `Available commands:
/help                           - Show this help
/quit                          - Exit interactive mode
/user create <id> <name>       - Create new user
/chat list                     - List your chats
/chat create <id> <u1>,<u2>    - Create new chat
/chat use <chat_id>            - Enter chat mode
/msg send <chat_id> <text>     - Send message to chat
/msg list <chat_id> [limit]    - List messages from chat

In chat mode (after /chat use):
<text>                         - Send message to current chat
/history [limit]               - Show message history
/members                       - Show chat members
/leave                         - Exit chat mode`

	s.display.ShowMessage(help)
}

func (s *InteractiveSession) handleChatCommand(args []string) {
	if len(args) == 0 {
		s.display.ShowError("Usage: /chat <list|create|use>")
		return
	}

	switch args[0] {
	case "list":
		s.listChats()
	case "create":
		if len(args) < 3 {
			s.display.ShowError("Usage: /chat create <chat_id> <u1>,<u2>,...")
			return
		}
		s.createChat(args[1], args[2])
	case "use":
		if len(args) < 2 {
			s.display.ShowError("Usage: /chat use <chat_id>")
			return
		}
		s.useChat(args[1])
	default:
		s.display.ShowError("Unknown chat command: " + args[0])
	}
}

func (s *InteractiveSession) handleUserCommand(args []string) {
	if len(args) == 0 {
		s.display.ShowError("Usage: /user <create>")
		return
	}

	switch args[0] {
	case "create":
		if len(args) < 3 {
			s.display.ShowError("Usage: /user create <id> <name>")
			return
		}
		s.createUser(args[1], strings.Join(args[2:], " "))
	default:
		s.display.ShowError("Unknown user command: " + args[0])
	}
}

func (s *InteractiveSession) handleMsgCommand(args []string) {
	if len(args) == 0 {
		s.display.ShowError("Usage: /msg <send|list>")
		return
	}

	switch args[0] {
	case "send":
		if len(args) < 3 {
			s.display.ShowError("Usage: /msg send <chat_id> <text>")
			return
		}
		s.sendMessage(args[1], strings.Join(args[2:], " "))
	case "list":
		if len(args) < 2 {
			s.display.ShowError("Usage: /msg list <chat_id> [limit]")
			return
		}
		limit := 10
		if len(args) > 2 {
			if l, err := strconv.Atoi(args[2]); err == nil {
				limit = l
			}
		}
		s.listMessages(args[1], limit)
	default:
		s.display.ShowError("Unknown msg command: " + args[0])
	}
}

func (s *InteractiveSession) handleHistoryCommand(args []string) {
	if s.currentChat == "" {
		s.display.ShowError("No active chat")
		return
	}
	limit := 10
	if len(args) > 0 {
		if l, err := strconv.Atoi(args[0]); err == nil {
			limit = l
		}
	}
	s.listMessages(s.currentChat, limit)
}

func (s *InteractiveSession) handleMembersCommand() {
	if s.currentChat == "" {
		s.display.ShowError("No active chat")
		return
	}
	s.showChatMembers(s.currentChat)
}

func (s *InteractiveSession) handleLeaveCommand() {
	if s.currentChat == "" {
		s.display.ShowError("No active chat")
		return
	}

	// Unsubscribe from current chat
	unsubMsg := WSMessage{
		Type:   "unsubscribe",
		ChatID: s.currentChat,
	}
	if err := s.conn.WriteJSON(unsubMsg); err != nil {
		s.display.ShowError("Failed to unsubscribe from chat")
	}

	s.display.ShowMessage(fmt.Sprintf("Left chat: %s", s.currentChat))
	s.currentChat = ""
}

func (s *InteractiveSession) listChats() {
	url := fmt.Sprintf("http://%s/chats?user_id=%s", s.serverAddr, s.userID)
	resp, err := http.Get(url)
	if err != nil {
		s.display.ShowError("Failed to fetch chats")
		return
	}
	defer resp.Body.Close()

	var chats []domain.Chat
	if err := json.NewDecoder(resp.Body).Decode(&chats); err != nil {
		s.display.ShowError("Failed to parse chats")
		return
	}

	if len(chats) == 0 {
		s.display.ShowMessage("No chats found")
		return
	}

	s.display.ShowMessage("Your chats:")
	for _, chat := range chats {
		members := strings.Join(chat.Members, ", ")
		s.display.ShowMessage(fmt.Sprintf("  %s: [%s]", chat.ID, members))
	}
}

func (s *InteractiveSession) createChat(chatID, membersStr string) {
	members := strings.Split(membersStr, ",")
	for i, m := range members {
		members[i] = strings.TrimSpace(m)
	}

	postJSON(fmt.Sprintf("http://%s", s.serverAddr), "/chats", map[string]interface{}{
		"id":      chatID,
		"members": members,
	})
	s.display.ShowMessage(fmt.Sprintf("Chat %s created", chatID))
}

func (s *InteractiveSession) useChat(chatID string) {
	// Unsubscribe from previous chat if any
	if s.currentChat != "" {
		unsubMsg := WSMessage{
			Type:   "unsubscribe",
			ChatID: s.currentChat,
		}
		s.conn.WriteJSON(unsubMsg)
	}

	s.currentChat = chatID
	s.display.ShowMessage(fmt.Sprintf("Entered chat: %s", chatID))

	// Subscribe to new chat
	subMsg := WSMessage{
		Type:   "subscribe",
		ChatID: chatID,
	}
	if err := s.conn.WriteJSON(subMsg); err != nil {
		s.display.ShowError("Failed to subscribe to chat")
		return
	}

	s.listMessages(chatID, 5)
}

func (s *InteractiveSession) createUser(userID, name string) {
	postJSON(fmt.Sprintf("http://%s", s.serverAddr), "/users", map[string]string{
		"id":   userID,
		"name": name,
	})
	s.display.ShowMessage(fmt.Sprintf("User %s created", userID))
}

func (s *InteractiveSession) sendMessage(chatID, text string) {
	msg := WSMessage{
		Type:   "message",
		ChatID: chatID,
		Text:   text,
	}

	if err := s.conn.WriteJSON(msg); err != nil {
		s.display.ShowError("Failed to send message")
	}
}

func (s *InteractiveSession) listMessages(chatID string, limit int) {
	url := fmt.Sprintf("http://%s/messages?user_id=%s&chat_id=%s", s.serverAddr, s.userID, chatID)
	resp, err := http.Get(url)
	if err != nil {
		s.display.ShowError("Failed to fetch messages")
		return
	}
	defer resp.Body.Close()

	var msgs []domain.Message
	if err := json.NewDecoder(resp.Body).Decode(&msgs); err != nil {
		s.display.ShowError("Failed to parse messages")
		return
	}

	if len(msgs) == 0 {
		s.display.ShowMessage("No messages found")
		return
	}

	start := 0
	if len(msgs) > limit {
		start = len(msgs) - limit
	}

	for _, msg := range msgs[start:] {
		s.display.ShowMessage(fmt.Sprintf("[%s] %s: %s", 
			msg.CreatedAt.Format("15:04:05"), msg.From, msg.Text))
	}
}

func (s *InteractiveSession) showChatMembers(chatID string) {
	url := fmt.Sprintf("http://%s/chats/%s", s.serverAddr, chatID)
	resp, err := http.Get(url)
	if err != nil {
		s.display.ShowError("Failed to fetch chat info")
		return
	}
	defer resp.Body.Close()

	var chat domain.Chat
	if err := json.NewDecoder(resp.Body).Decode(&chat); err != nil {
		s.display.ShowError("Failed to parse chat info")
		return
	}

	s.display.ShowMessage(fmt.Sprintf("Members of %s: %s", chatID, strings.Join(chat.Members, ", ")))
}

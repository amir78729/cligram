package api

import (
	"cligram/cmd/server/types"
	"cligram/internal/app"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	Service *app.ChatService
}

func (s *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req types.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateUserHandler decode error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("CreateUserHandler: creating user %s", req.ID)
	if err := s.Service.CreateUser(req.ID, req.Name); err != nil {
		log.Printf("CreateUserHandler error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("CreateUserHandler: user %s created successfully", req.ID)
}

func (s *Server) CreateChatHandler(w http.ResponseWriter, r *http.Request) {
	var req types.CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateChatHandler decode error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("CreateChatHandler: creating chat %s with members %v", req.ID, req.Members)
	if err := s.Service.CreateChat(req.ID, req.Members); err != nil {
		log.Printf("CreateChatHandler error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("CreateChatHandler: chat %s created successfully", req.ID)
}

func (s *Server) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var req types.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("SendMessageHandler decode error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("SendMessageHandler: sending message from %s to chat %s", req.From, req.ChatID)
	if err := s.Service.SendMessage(req.From, req.ChatID, req.Text); err != nil {
		log.Printf("SendMessageHandler error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Printf("SendMessageHandler: message sent successfully from %s to chat %s", req.From, req.ChatID)
}

func (s *Server) ListMessagesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	chatID := r.URL.Query().Get("chat_id")

	if userID == "" || chatID == "" {
		log.Printf("ListMessagesHandler missing parameters: user_id=%s chat_id=%s", userID, chatID)
		http.Error(w, "user_id and chat_id are required", http.StatusBadRequest)
		return
	}

	log.Printf("ListMessagesHandler: listing messages for user %s in chat %s", userID, chatID)
	msgs, err := s.Service.ListMessages(userID, chatID)
	if err != nil {
		log.Printf("ListMessagesHandler error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(msgs); err != nil {
		log.Printf("ListMessagesHandler encode error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("ListMessagesHandler: returned %d messages for user %s in chat %s", len(msgs), userID, chatID)
}

func (s *Server) ListChatsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Printf("ListChatsHandler missing user_id parameter")
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	log.Printf("ListChatsHandler: listing chats for user %s", userID)
	chats, err := s.Service.ListUserChats(userID)
	if err != nil {
		log.Printf("ListChatsHandler error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(chats); err != nil {
		log.Printf("ListChatsHandler encode error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("ListChatsHandler: returned %d chats for user %s", len(chats), userID)
}

func (s *Server) GetChatHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatID := vars["id"]
	
	if chatID == "" {
		log.Printf("GetChatHandler missing chat_id in path")
		http.Error(w, "chat_id is required", http.StatusBadRequest)
		return
	}

	log.Printf("GetChatHandler: getting chat %s", chatID)
	chat, err := s.Service.GetChatByID(chatID)
	if err != nil {
		log.Printf("GetChatHandler error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(chat); err != nil {
		log.Printf("GetChatHandler encode error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("GetChatHandler: returned chat %s", chatID)
}

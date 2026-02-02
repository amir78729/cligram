package domain

import "time"

type User struct {
	ID   string
	Name string
}

type Message struct {
	ID        string
	From      string
	ChatID    string
	Text      string
	CreatedAt time.Time
}

type Chat struct {
	ID      string
	Members []string
}

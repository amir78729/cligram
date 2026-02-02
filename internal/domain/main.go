package domain

import "time"

type User struct {
	ID   string
	Name string
}

type Message struct {
	ID        string    `json:"id" bson:"_id"`
	From      string    `json:"from" bson:"from"`
	ChatID    string    `json:"chat_id" bson:"chat_id"`
	Text      string    `json:"text" bson:"text"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type Chat struct {
	ID      string
	Members []string
}

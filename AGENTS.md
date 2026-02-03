# AGENTS.md - Cligram Project Context

## Project Overview

**Cligram** is a CLI-based messenger application built in Go for learning system design and Go programming. The project implements a client-server architecture with real-time messaging capabilities.

## Architecture

### Core Components

- **Domain Layer**: Entities (User, Chat, Message) and business rules
- **Repository Layer**: Data access abstractions
- **Service Layer**: Business logic (ChatService)
- **API Layer**: REST endpoints + WebSocket handlers
- **CLI Layer**: Command-line interface and interactive chat

### Technology Stack

- **Language**: Go 1.25.5
- **Database**: MongoDB with official driver
- **WebSocket**: Gorilla WebSocket
- **HTTP Router**: Gorilla Mux
- **CLI**: Custom implementation with interactive mode
- **Styling**: Gookit Color for terminal output

## Project Structure

```
cligram/
├── cmd/
│   ├── cligram/          # CLI client entry point
│   └── server/           # Server entry point and API handlers
├── internal/
│   ├── domain/           # Core entities and business rules
│   │   └── repository/   # Repository interfaces
│   ├── db/              # MongoDB implementations
│   ├── app/             # Service layer (ChatService)
│   ├── cli/             # CLI commands and interactive mode
│   └── theme/           # UI styling and tokens
```

## Key Entities

### User

```go
type User struct {
    ID   string
    Name string
}
```

### Chat

```go
type Chat struct {
    ID      string
    Members []string
}
```

### Message

```go
type Message struct {
    ID        string    `json:"id" bson:"_id"`
    From      string    `json:"from" bson:"from"`
    ChatID    string    `json:"chat_id" bson:"chat_id"`
    Text      string    `json:"text" bson:"text"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
```

## Current Features

### CLI Commands

- `cligram user <create|get> [args]` - User management
- `cligram chat <create|get> [args]` - Chat management
- `cligram msg <send|list> [args]` - Message operations
- `cligram interactive <user_id> <server_addr>` - Command-based interactive chat

### Interactive Commands

- `/help` - Show available commands
- `/quit` - Exit interactive mode
- `/user create <id> <name>` - Create new user
- `/chat list` - List user's chats with members
- `/chat create <id> <u1>,<u2>` - Create new chat
- `/chat use <chat_id>` - Enter chat mode (stateful)
- `/msg send <chat_id> <text>` - Send message to specific chat
- `/msg list <chat_id> [limit]` - List recent messages

### Chat Mode Commands (after `/chat use <chat_id>`)

- `<text>` - Send message to current chat
- `/history [limit]` - Show message history
- `/members` - Show chat members
- `/leave` - Exit chat mode

### Server Endpoints

- `POST /users` - Create user
- `POST /chats` - Create chat
- `GET /chats?user_id=<id>` - List user's chats
- `GET /chats/{id}` - Get chat details
- `POST /messages` - Send message
- `GET /messages?user_id=<id>&chat_id=<id>` - List messages
- `GET /ws?user_id=<id>` - WebSocket connection

### Real-time Features

- WebSocket-based messaging with dynamic subscription control
- Manual chat subscription via `/chat use <chat_id>` command
- Manual unsubscription via `/leave` command
- Message broadcasting only to subscribed chat members
- Command-based interactive CLI with live message updates
- Stateful chat mode for focused conversations
- Clean subscription management (no automatic subscriptions)

## Business Logic (ChatService)

- User creation and validation
- Chat creation with member validation
- Message sending with authorization checks
- Message listing with access control
- User chat listing

## Database Layer

- MongoDB collections: users, chats, messages
- Repository pattern implementation
- Proper error handling and domain error mapping

## Current Limitations & Growth Areas

1. **Authentication**: No user authentication system
2. **Authorization**: Basic member-based access control only
3. **Persistence**: No message history pagination
4. **Configuration**: Hardcoded server settings
5. **Testing**: No test coverage
6. **Logging**: Basic logging implementation
7. **Error Handling**: Could be more comprehensive
8. **Validation**: Limited input validation
9. **Security**: No rate limiting or input sanitization
10. **Scalability**: Single-server architecture

## Development Guidelines

- Follow clean architecture principles
- Maintain separation of concerns
- Use dependency injection
- Keep domain logic pure
- Implement proper error handling
- Write minimal, focused code
- Use adapter pattern for replaceable components (DisplayManager)
- Update this AGENTS.md file after significant changes

## Next Development Priorities

1. Authentication system
2. Enhanced error handling
3. Configuration management
4. Testing framework
5. Message pagination
6. User presence indicators
7. Chat room features
8. File sharing capabilities

---

_Last Updated: 2026-02-03_
_Version: Initial implementation with WebSocket support and subscription control_

- Use adapter pattern for replaceable components (DisplayManager)
- Update this AGENTS.md file after significant changes

## Next Development Priorities

1. Authentication system
2. Enhanced error handling
3. Configuration management
4. Testing framework
5. Message pagination
6. User presence indicators
7. Chat room features
8. File sharing capabilities

---

_Last Updated: 2026-02-03_
_Version: Initial implementation with WebSocket support and subscription control_
_Version: Initial implementation with WebSocket support and subscription control_

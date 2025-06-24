# Bagchal Game Backend

This is a Golang backend server for the Bagchal game, supporting both player vs AI and player vs player modes.

## Features

- RESTful API for game management
- WebSocket support for real-time game updates
- AI opponent implementation
- Game state persistence
- Move validation
- Game history tracking

## Project Structure

```
.
├── cmd/
│   └── server/           # Main application entry point
├── internal/
│   ├── ai/              # AI logic implementation
│   ├── api/             # API handlers and routes
│   ├── domain/          # Core game logic and models
│   ├── config/          # Configuration management
│   └── websocket/       # WebSocket handling
└── pkg/                 # Shared packages
```

## Setup

1. Install Go 1.21 or later
2. Clone the repository
3. Run `go mod init github.com/yourusername/BagchalGolang`
4. Run `go mod tidy` to install dependencies
5. Run `go run cmd/server/main.go` to start the server

## API Documentation

POST /api/games - 創建新遊戲
GET /api/games/:id - 獲取遊戲狀態
POST /api/games/:id/moves - 執行移動
GET /api/games/player/:playerID - 獲取玩家的遊戲列表
DELETE /api/games/:id - 刪除遊戲

## Game Rules

Bagchal is a traditional board game from Nepal. Here are the basic rules:

- Played on a 5x5 board
- Two players: Tigers (4 pieces) and Goats (20 pieces)
- Tigers can move to any adjacent intersection
- Tigers can capture goats by jumping over them
- Goats can only move to adjacent intersections
- Goats win by blocking all tiger moves
- Tigers win by capturing 5 goats 
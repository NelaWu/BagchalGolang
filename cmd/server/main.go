package main

import (
	"github.com/gin-contrib/cors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nelawu/BagchalGolang/internal/ai"
	"github.com/nelawu/BagchalGolang/internal/api/handler"
	"github.com/nelawu/BagchalGolang/internal/domain/game"
)

// MemoryGameRepository 內存遊戲存儲實現
type MemoryGameRepository struct {
	games map[string]*game.Game
	mu    sync.RWMutex
}

func NewMemoryGameRepository() *MemoryGameRepository {
	return &MemoryGameRepository{
		games: make(map[string]*game.Game),
	}
}

func (r *MemoryGameRepository) Save(game *game.Game) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.games[game.ID] = game
	return nil
}

func (r *MemoryGameRepository) GetByID(id string) (*game.Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if game, exists := r.games[id]; exists {
		return game, nil
	}
	return nil, game.ErrGameNotFound
}

func (r *MemoryGameRepository) List(playerID string) ([]*game.Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var playerGames []*game.Game
	for _, g := range r.games {
		if g.PlayerID == playerID {
			playerGames = append(playerGames, g)
		}
	}
	return playerGames, nil
}

func (r *MemoryGameRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.games, id)
	return nil
}

func main() {
	// 創建路由
	router := gin.Default()

	// 配置 CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // 允許的來源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 初始化依賴
	gameRepo := NewMemoryGameRepository()
	aiEngine := ai.NewEngine(2) // 默認中等難度
	gameService := game.NewGameService(gameRepo, aiEngine)
	gameHandler := handler.NewGameHandler(gameService)

	// 註冊路由
	gameHandler.RegisterRoutes(router)

	// 啟動服務器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("啟動服務器失敗: %v", err)
	}
}

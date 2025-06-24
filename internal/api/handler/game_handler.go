package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nelawu/BagchalGolang/internal/domain/game"
)

type GameHandler struct {
	gameService *game.GameService
}

func NewGameHandler(gameService *game.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// RegisterRoutes 註冊路由
func (h *GameHandler) RegisterRoutes(router *gin.Engine) {
	gameGroup := router.Group("/api/games")
	{
		gameGroup.POST("", h.createGame)
		gameGroup.GET("/:id", h.getGame)
		gameGroup.POST("/:id/moves", h.makeMove)
		gameGroup.GET("/player/:playerID", h.listPlayerGames)
		gameGroup.DELETE("/:id", h.deleteGame)
	}
}

// CreateGameRequest 創建遊戲請求
type CreateGameRequest struct {
	PlayerID  string `json:"playerId"`
	IsAIGame  bool   `json:"isAIGame"`
	AILevel   int    `json:"aiLevel"`
}

// createGame 創建新遊戲
func (h *GameHandler) createGame(c *gin.Context) {
	var req CreateGameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求數據"})
		return
	}

	// 驗證AI難度級別
	if req.IsAIGame && (req.AILevel < 1 || req.AILevel > 3) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的AI難度級別（應為1-3）"})
		return
	}

	newGame, err := h.gameService.CreateGame(req.PlayerID, req.IsAIGame, req.AILevel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "創建遊戲失敗"})
		return
	}

	c.JSON(http.StatusCreated, newGame)
}

// getGame 獲取遊戲
func (h *GameHandler) getGame(c *gin.Context) {
	gameID := c.Param("id")
	game, err := h.gameService.GetGame(gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "遊戲不存在"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// MakeMoveRequest 移動請求
type MakeMoveRequest struct {
	From      game.Position `json:"from"`
	To        game.Position `json:"to"`
	PieceType game.PieceType `json:"pieceType"`
}

// makeMove 執行移動
func (h *GameHandler) makeMove(c *gin.Context) {
	gameID := c.Param("id")
	var req MakeMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求數據"})
		return
	}

	move := game.Move{
		From:      req.From,
		To:        req.To,
		PieceType: req.PieceType,
	}

	updatedGame, err := h.gameService.MakeMove(gameID, move)
	if err != nil {
		switch err {
		case game.ErrGameNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "遊戲不存在"})
		case game.ErrInvalidMove:
			c.JSON(http.StatusBadRequest, gin.H{"error": "無效的移動"})
		case game.ErrGameOver:
			c.JSON(http.StatusBadRequest, gin.H{"error": "遊戲已結束"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "執行移動失敗"})
		}
		return
	}

	c.JSON(http.StatusOK, updatedGame)
}

// listPlayerGames 列出玩家的所有遊戲
func (h *GameHandler) listPlayerGames(c *gin.Context) {
	playerID := c.Param("playerID")
	games, err := h.gameService.ListPlayerGames(playerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "獲取遊戲列表失敗"})
		return
	}

	c.JSON(http.StatusOK, games)
}

// deleteGame 刪除遊戲
func (h *GameHandler) deleteGame(c *gin.Context) {
	gameID := c.Param("id")
	err := h.gameService.DeleteGame(gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除遊戲失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "遊戲已刪除"})
} 
package game

import (
	"time"
)

// 棋盤大小常數
const (
	BoardSize = 5
	MaxGoats  = 20
	MaxTigers = 4
)

// PieceType 表示棋子類型
type PieceType int

const (
	Empty PieceType = iota
	Tiger
	Goat
)

// Position 表示棋盤上的位置
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Move 表示一步棋
type Move struct {
	From     Position  `json:"from"`
	To       Position  `json:"to"`
	Capture  *Position `json:"capture,omitempty"` // 如果是虎吃羊，這裡記錄被吃的羊的位置
	PieceType PieceType `json:"pieceType"`
}

// GameState 表示遊戲狀態
type GameState struct {
	Board         [BoardSize][BoardSize]PieceType `json:"board"`
	GoatsInHand   int                            `json:"goatsInHand"`   // 還未放置的羊數量
	CapturedGoats int                            `json:"capturedGoats"` // 被吃掉的羊數量
	CurrentTurn   PieceType                      `json:"currentTurn"`   // 當前回合：虎或羊
	IsGameOver    bool                           `json:"isGameOver"`
	Winner        PieceType                      `json:"winner"`
	LastMove      *Move                          `json:"lastMove"`
}

// Game 表示一局遊戲
type Game struct {
	ID        string     `json:"id"`
	State     GameState  `json:"state"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	PlayerID  string     `json:"playerId"`  // 玩家ID
	IsAIGame  bool       `json:"isAIGame"`  // 是否是AI對戰
	AILevel   int        `json:"aiLevel"`   // AI難度等級
}

// NewGame 創建一個新遊戲
func NewGame(playerID string, isAIGame bool, aiLevel int) *Game {
	game := &Game{
		ID:        generateGameID(), // 需要實現此函數
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		PlayerID:  playerID,
		IsAIGame:  isAIGame,
		AILevel:   aiLevel,
	}

	// 初始化遊戲狀態
	game.State = GameState{
		GoatsInHand:   MaxGoats,
		CapturedGoats: 0,
		CurrentTurn:   Goat, // 羊先手
		IsGameOver:    false,
	}

	// 放置初始虎子
	game.State.Board[0][0] = Tiger
	game.State.Board[0][BoardSize-1] = Tiger
	game.State.Board[BoardSize-1][0] = Tiger
	game.State.Board[BoardSize-1][BoardSize-1] = Tiger

	return game
}

// IsValidMove 檢查移動是否合法
func (g *Game) IsValidMove(move Move) bool {
	// TODO: 實現移動驗證邏輯
	return true
}

// MakeMove 執行一步移動
func (g *Game) MakeMove(move Move) error {
	// TODO: 實現移動邏輯
	return nil
}

// generateGameID 生成遊戲ID
func generateGameID() string {
	// TODO: 實現遊戲ID生成邏輯
	return "game_" + time.Now().Format("20060102150405")
} 
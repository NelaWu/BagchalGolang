package game

// GameRepository 定義遊戲資料存儲介面
type GameRepository interface {
	// Save 保存遊戲狀態
	Save(game *Game) error

	// GetByID 根據ID獲取遊戲
	GetByID(id string) (*Game, error)

	// List 列出玩家的所有遊戲
	List(playerID string) ([]*Game, error)

	// Delete 刪除遊戲
	Delete(id string) error
} 
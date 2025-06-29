package game

import (
	"errors"
	"log"
	"math"
)

var (
	ErrInvalidMove    = errors.New("invalid move")
	ErrGameNotFound   = errors.New("game not found")
	ErrNotPlayersTurn = errors.New("not player's turn")
	ErrGameOver       = errors.New("game is already over")
)

type GameService struct {
	repository GameRepository
	aiEngine   AIEngine
}

type AIEngine interface {
	CalculateNextMove(game *Game) (*Move, error)
}

func NewGameService(repository GameRepository, aiEngine AIEngine) *GameService {
	return &GameService{
		repository: repository,
		aiEngine:   aiEngine,
	}
}

// CreateGame 創建新遊戲
func (s *GameService) CreateGame(playerID string, isAIGame bool, aiLevel int) (*Game, error) {
	game := NewGame(playerID, isAIGame, aiLevel)
	err := s.repository.Save(game)
	if err != nil {
		return nil, err
	}
	return game, nil
}

// MakeMove 執行移動並處理遊戲邏輯
func (s *GameService) MakeMove(gameID string, move Move) (*Game, error) {
	game, err := s.repository.GetByID(gameID)
	if err != nil {
		return nil, ErrGameNotFound
	}

	if game.State.IsGameOver {
		return nil, ErrGameOver
	}

	if !s.IsValidMove(game, move) {
		return nil, ErrInvalidMove
	}

	// 執行移動
	err = s.executeMove(game, move)
	if err != nil {
		return nil, err
	}

	// 檢查遊戲是否結束
	s.checkGameOver(game)

	// 如果是AI遊戲且遊戲未結束，執行AI移動
	if game.IsAIGame && !game.State.IsGameOver {
		aiMove, err := s.aiEngine.CalculateNextMove(game)
		if err != nil {
			return nil, err
		}
		err = s.executeMove(game, *aiMove)
		if err != nil {
			return nil, err
		}
		s.checkGameOver(game)
	}

	// 保存遊戲狀態
	err = s.repository.Save(game)
	if err != nil {
		return nil, err
	}

	return game, nil
}

// IsValidMove 檢查移動是否合法
func (s *GameService) IsValidMove(game *Game, move Move) bool {
	// 檢查是否是玩家的回合
	if game.State.CurrentTurn != move.PieceType {
		return false
	}

	// 檢查目標位置是否為空
	if game.State.Board[move.To.Y][move.To.X] != Empty {
		return false
	}

	// 如果是放置羊的階段
	if move.PieceType == Goat && game.State.GoatsInHand > 0 {
		return move.From.X == move.To.X && move.From.Y == move.To.Y
	}

	// 對於非放置階段的移動，檢查起點是否有正確的棋子
	if game.State.Board[move.From.Y][move.From.X] != move.PieceType {
		return false
	}

	// 檢查起點是否有正確的棋子
	if game.State.Board[move.From.Y][move.From.X] != move.PieceType {
		return false
	}

	// 檢查目標位置是否為空
	if game.State.Board[move.To.Y][move.To.X] != Empty {
		return false
	}

	// 如果是放置羊的階段
	if move.PieceType == Goat && game.State.GoatsInHand > 0 {
		return move.From.X == move.To.X && move.From.Y == move.To.Y
	}

	// 檢查移動距離
	dx := math.Abs(float64(move.To.X - move.From.X))
	dy := math.Abs(float64(move.To.Y - move.From.Y))

	// 正常移動：只能移動到相鄰的點
	if dx <= 1 && dy <= 1 {
		return true
	}

	// 虎吃羊：檢查是否是有效的跳躍
	if move.PieceType == Tiger && dx <= 2 && dy <= 2 {
		// 計算中間位置
		midX := (move.From.X + move.To.X) / 2
		midY := (move.From.Y + move.To.Y) / 2

		// 檢查中間是否有羊
		if game.State.Board[midY][midX] == Goat {
			move.Capture = &Position{X: midX, Y: midY}
			return true
		}
	}

	return false
}

// executeMove 執行移動
func (s *GameService) executeMove(game *Game, move Move) error {
	// 如果是放置羊的階段
	if move.PieceType == Goat && game.State.GoatsInHand > 0 {
		game.State.Board[move.To.Y][move.To.X] = Goat
		game.State.GoatsInHand--
	} else {
		// 移動棋子
		game.State.Board[move.From.Y][move.From.X] = Empty
		game.State.Board[move.To.Y][move.To.X] = move.PieceType

		// 處理吃子
		isValid := (move.To.X+move.From.X)%2 == 0 && (move.To.Y+move.From.Y)%2 == 0
		if isValid == true && game.State.Board[(move.To.Y+move.From.Y)/2][(move.To.X+move.From.X)/2] == Goat {
			game.State.Board[(move.To.Y+move.From.Y)/2][(move.To.X+move.From.X)/2] = Empty
			game.State.CapturedGoats++
		}
	}

	// 更新最後一步
	game.State.LastMove = &move

	// 切換回合
	if game.State.CurrentTurn == Tiger {
		game.State.CurrentTurn = Goat
	} else {
		game.State.CurrentTurn = Tiger
	}
	log.Printf("executeMove2：%d", game.State.CurrentTurn)

	return nil
}

// checkGameOver 檢查遊戲是否結束
func (s *GameService) checkGameOver(game *Game) {
	// 檢查虎是否獲勝（吃掉5隻羊）
	if game.State.CapturedGoats >= 5 {
		log.Printf("檢查虎是否獲勝（吃掉5隻羊）")
		game.State.IsGameOver = true
		game.State.Winner = Tiger
		return
	}

	// 檢查羊是否獲勝（虎無法移動）
	if s.isTigerTrapped(game) || game.State.GoatsInHand == 0 {
		log.Printf("檢查羊是否獲勝（虎無法移動）")
		game.State.IsGameOver = true
		game.State.Winner = Goat
		return
	}
}

// isTigerTrapped 檢查虎是否無法移動
func (s *GameService) isTigerTrapped(game *Game) bool {
	// 對每隻虎檢查是否有可行的移動
	for y := 0; y < BoardSize; y++ {
		for x := 0; x < BoardSize; x++ {
			if game.State.Board[y][x] == Tiger {
				if s.hasTigerMoves(game, x, y) {
					return false
				}
			}
		}
	}
	return true
}

// hasTigerMoves 檢查特定位置的虎是否有可行的移動
func (s *GameService) hasTigerMoves(game *Game, x, y int) bool {
	// 檢查所有可能的移動方向
	log.Printf("檢查虎子在位置(%d,%d)的可能移動", x, y)
	directions := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, d := range directions {
		newX, newY := x+d.dx, y+d.dy
		// 檢查普通移動
		if newX >= 0 && newX < BoardSize && newY >= 0 && newY < BoardSize {
			log.Printf("檢查普通移動到(%d,%d)", newX, newY)
			if game.State.Board[newY][newX] == Empty {
				log.Printf("找到有效的普通移動")
				return true
			}
		}

		// 檢查跳躍移動（吃子）
		jumpX, jumpY := x+2*d.dx, y+2*d.dy
		if jumpX >= 0 && jumpX < BoardSize && jumpY >= 0 && jumpY < BoardSize {
			if game.State.Board[jumpY][jumpX] == Empty {
				return true
			}
		}
	}

	log.Printf("該虎子沒有可用的移動")
	return false
}

// GetGame 根據ID獲取遊戲
func (s *GameService) GetGame(id string) (*Game, error) {
	return s.repository.GetByID(id)
}

// ListPlayerGames 獲取玩家的所有遊戲
func (s *GameService) ListPlayerGames(playerID string) ([]*Game, error) {
	return s.repository.List(playerID)
}

// DeleteGame 刪除遊戲
func (s *GameService) DeleteGame(id string) error {
	return s.repository.Delete(id)
}

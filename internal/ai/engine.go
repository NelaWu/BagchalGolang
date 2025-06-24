package ai

import (
	"math/rand"
	"time"

	"github.com/nelawu/BagchalGolang/internal/domain/game"
)

type Engine struct {
	difficulty int
}

func NewEngine(difficulty int) *Engine {
	return &Engine{
		difficulty: difficulty,
	}
}

// CalculateNextMove 計算AI的下一步移動
func (e *Engine) CalculateNextMove(g *game.Game) (*game.Move, error) {
	// TODO: 實現更智能的AI邏輯
	// 目前僅實現一個簡單的隨機移動策略

	// 初始化隨機數生成器
	rand.Seed(time.Now().UnixNano())

	// 確定當前是虎方還是羊方
	isGoat := g.State.CurrentTurn == game.Goat

	var validMoves []game.Move

	// 如果是羊方且還有羊可以放置
	if isGoat && g.State.GoatsInHand > 0 {
		// 收集所有可能的放置位置
		for y := 0; y < game.BoardSize; y++ {
			for x := 0; x < game.BoardSize; x++ {
				if g.State.Board[y][x] == game.Empty {
					move := game.Move{
						From:      game.Position{X: x, Y: y},
						To:        game.Position{X: x, Y: y},
						PieceType: game.Goat,
					}
					validMoves = append(validMoves, move)
				}
			}
		}
	} else {
		// 收集所有可能的移動
		pieceType := game.Tiger
		if isGoat {
			pieceType = game.Goat
		}

		// 遍歷棋盤尋找己方棋子
		for y := 0; y < game.BoardSize; y++ {
			for x := 0; x < game.BoardSize; x++ {
				if g.State.Board[y][x] == pieceType {
					// 檢查所有可能的移動方向
					directions := []struct{ dx, dy int }{
						{-1, -1}, {-1, 0}, {-1, 1},
						{0, -1}, {0, 1},
						{1, -1}, {1, 0}, {1, 1},
					}

					for _, d := range directions {
						// 檢查普通移動
						newX, newY := x+d.dx, y+d.dy
						if isValidPosition(newX, newY) {
							move := game.Move{
								From:      game.Position{X: x, Y: y},
								To:        game.Position{X: newX, Y: newY},
								PieceType: pieceType,
							}
							if isValidMove(g, move) {
								validMoves = append(validMoves, move)
							}
						}

						// 如果是虎，還要檢查吃子移動
						if pieceType == game.Tiger {
							jumpX, jumpY := x+2*d.dx, y+2*d.dy
							if isValidPosition(jumpX, jumpY) {
								move := game.Move{
									From:      game.Position{X: x, Y: y},
									To:        game.Position{X: jumpX, Y: jumpY},
									PieceType: game.Tiger,
								}
								if isValidMove(g, move) {
									// 優先考慮吃子移動
									validMoves = append(validMoves, move)
								}
							}
						}
					}
				}
			}
		}
	}

	// 如果沒有有效的移動，返回錯誤
	if len(validMoves) == 0 {
		return nil, nil
	}

	// 根據難度選擇移動
	var selectedMove game.Move
	switch e.difficulty {
	case 1: // 簡單：完全隨機
		selectedMove = validMoves[rand.Intn(len(validMoves))]
	case 2: // 中等：優先選擇吃子移動
		// 將吃子移動放在前面
		var captureMoves []game.Move
		for _, move := range validMoves {
			if move.Capture != nil {
				captureMoves = append(captureMoves, move)
			}
		}
		if len(captureMoves) > 0 {
			selectedMove = captureMoves[rand.Intn(len(captureMoves))]
		} else {
			selectedMove = validMoves[rand.Intn(len(validMoves))]
		}
	case 3: // 困難：TODO: 實現更智能的策略
		// 目前與中等難度相同
		selectedMove = validMoves[rand.Intn(len(validMoves))]
	default:
		selectedMove = validMoves[rand.Intn(len(validMoves))]
	}

	return &selectedMove, nil
}

// isValidPosition 檢查位置是否在棋盤範圍內
func isValidPosition(x, y int) bool {
	return x >= 0 && x < game.BoardSize && y >= 0 && y < game.BoardSize
}

// isValidMove 檢查移動是否合法
func isValidMove(g *game.Game, move game.Move) bool {
	// 檢查目標位置是否為空
	if g.State.Board[move.To.Y][move.To.X] != game.Empty {
		return false
	}

	// 如果是放置羊的階段
	if move.PieceType == game.Goat && g.State.GoatsInHand > 0 {
		return move.From.X == move.To.X && move.From.Y == move.To.Y
	}

	// 計算移動距離
	dx := abs(move.To.X - move.From.X)
	dy := abs(move.To.Y - move.From.Y)

	// 正常移動：只能移動到相鄰的點
	if dx <= 1 && dy <= 1 {
		return true
	}

	// 虎吃羊：檢查是否是有效的跳躍
	if move.PieceType == game.Tiger && dx <= 2 && dy <= 2 {
		// 計算中間位置
		midX := (move.From.X + move.To.X) / 2
		midY := (move.From.Y + move.To.Y) / 2

		// 檢查中間是否有羊
		if g.State.Board[midY][midX] == game.Goat {
			move.Capture = &game.Position{X: midX, Y: midY}
			return true
		}
	}

	return false
}

// abs 返回整數的絕對值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
} 
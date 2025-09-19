package service

import (
	"context"
	"crypto/rand"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog"
	"math/big"
	"narde/internal/domain"
	"time"
)

type GameService struct {
	gameStateCacher GameStateCacher
}

func NewGameService() *GameService {
	return &GameService{}
}

func (s *GameService) CreateGameState(ctx context.Context, player1, player2 uuid.UUID) (*domain.GameState, error) {
	logger := zerolog.Ctx(ctx)

	gameID, err := uuid.NewV4()
	if err != nil {
		logger.Error().Err(err).Send()
		return nil, err
	}

	gameState := &domain.GameState{
		GameID:        gameID,
		WhitePlayerID: player1,
		BlackPlayerID: player2,
		Board:         [24]int{15, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -15},
		Turn:          domain.WhiteTurnType,
		Dice:          rollDice(),
		LastMoveAt:    time.Now(),
	}

	//TODO save gameState in cache
	return gameState, nil
}

func (s *GameService) MakeMove(ctx context.Context, userID, gameID uuid.UUID, moves []domain.Move) (*domain.GameState, error) {
	logger := zerolog.Ctx(ctx)

	gameState, err := s.gameStateCacher.GetGameStateByID(ctx, gameID)
	if err != nil {
		logger.Error().Err(err).Send()
		return nil, err
	}

	availableMovesMap := map[int]int{
		gameState.Dice[0]: 1,
		gameState.Dice[1]: 1,
	}

	numberOfMoves := 2
	if gameState.Dice[0] == gameState.Dice[1] {
		numberOfMoves = 4
		availableMovesMap[gameState.Dice[0]] = 4
	}

	if len(moves) > numberOfMoves || len(moves) == 0 {
		return nil, domain.ErrInvalidMove
	}

	activePlayer := gameState.WhitePlayerID
	if gameState.Turn == domain.BlackTurnType {
		activePlayer = gameState.BlackPlayerID
	}

	if activePlayer != userID {
		return nil, domain.ErrInvalidMove
	}

	boardCopy := gameState.Board
	for _, move := range moves {
		if gameState.Turn == domain.WhiteTurnType {
			if boardCopy[move.From] <= 0 || boardCopy[move.To] < 0 {
				return nil, domain.ErrInvalidMove
			}
			dist := move.To - move.From
			if val, ok := availableMovesMap[dist]; !ok || val <= 0 {
				return nil, domain.ErrInvalidMove

			}
			piece := 1
			availableMovesMap[dist] -= 1
			boardCopy[move.From] -= piece
			boardCopy[move.To] += piece

		} else {
			if boardCopy[move.From] >= 0 || boardCopy[move.To] > 0 {
				return nil, domain.ErrInvalidMove
			}
			dist := -(move.To - move.From)
			if val, ok := availableMovesMap[dist]; !ok || val <= 0 {
				return nil, domain.ErrInvalidMove

			}
			piece := -1
			availableMovesMap[dist] -= 1
			boardCopy[move.From] -= piece
			boardCopy[move.To] += piece
		}

	}

	gameState.Board = boardCopy
	gameState.Dice = rollDice()
	gameState.Turn = changeTurn(gameState.Turn)
	gameState.LastMoveAt = time.Now()
	if err := s.gameStateCacher.SaveGameState(ctx, gameState); err != nil {
		logger.Error().Err(err).Send()
		return nil, err
	}

	return gameState, nil
}

func changeTurn(currentTurn domain.TurnType) domain.TurnType {
	if currentTurn == domain.WhiteTurnType {
		return domain.BlackTurnType
	}
	return domain.WhiteTurnType
}

func rollDice() [2]int {
	die1, _ := rand.Int(rand.Reader, big.NewInt(6))
	die2, _ := rand.Int(rand.Reader, big.NewInt(6))
	return [2]int{int(die1.Int64()) + 1, int(die2.Int64()) + 1}
}

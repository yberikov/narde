package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/rs/zerolog"
	"narde/internal/domain"
)

type (
	MatchmakingService interface {
		FindOpponent(ctx context.Context, playerID uuid.UUID) (uuid.UUID, uuid.UUID, error)
	}
	GameService interface {
		MakeMove(ctx context.Context, userID, gameID uuid.UUID, moves []domain.Move) (*domain.GameState, error)
		CreateGameState(ctx context.Context, player1, player2 uuid.UUID) (*domain.GameState, error)
	}
	Hub struct {
		Clients            map[uuid.UUID]*Client //userID
		Register           chan *Client
		Unregister         chan *Client
		Matchmaking        chan *Client
		GameMove           chan *MoveEvent
		matchmakingService MatchmakingService
		gameService        GameService
	}
)

func NewHub(
	matchmakingService MatchmakingService,
	gameService GameService,
) *Hub {
	return &Hub{
		Clients:            make(map[uuid.UUID]*Client),
		Register:           make(chan *Client),
		Unregister:         make(chan *Client),
		Matchmaking:        make(chan *Client),
		GameMove:           make(chan *MoveEvent),
		matchmakingService: matchmakingService,
		gameService:        gameService,
	}
}

func (h *Hub) Run() {
	ctx := context.Background()
	logger := zerolog.Ctx(ctx)

	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = client
			client.Send <- []byte("Hello")
			logger.Info().Msgf("🔌 New client connected %s", client.UserID.String())
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				logger.Info().Msgf("Client disconnected: %s", client.UserID.String())
			}

		case client := <-h.Matchmaking:
			player1, player2, err := h.matchmakingService.FindOpponent(ctx, client.UserID)
			if err != nil && !errors.Is(err, domain.ErrOpponentNotFound) {
				logger.Error().Err(err).Send()
				continue
			}

			if errors.Is(err, domain.ErrOpponentNotFound) {
				logger.Info().Msgf("Player %s is now waiting in the queue", client.UserID)
				waitMsg, _ := json.Marshal(Message{
					Type:    "matchmaking_wait",
					Payload: json.RawMessage("Wait"),
				})
				client.Send <- waitMsg
			}

			gameState, err := h.gameService.CreateGameState(ctx, player1, player2)

			logger.Info().Msgf("Match start", client.UserID)
			startMsg, _ := json.Marshal(Message{
				Type:    "match_start",
				Payload: json.RawMessage("Match Start"),
			})

			whitePlayerClient, err := h.findClientByID(gameState.WhitePlayerID)
			if err != nil {
				logger.Error().Err(err).Send()
				continue
			}

			blackPlayerClient, err := h.findClientByID(gameState.BlackPlayerID)
			if err != nil {
				logger.Error().Err(err).Send()
				continue
			}

			whitePlayerClient.GameID = gameState.GameID
			blackPlayerClient.GameID = gameState.GameID

			whitePlayerClient.Send <- startMsg
			blackPlayerClient.Send <- startMsg

		case moveEvent := <-h.GameMove:
			client := moveEvent.Client

			var payload MakeMovePayload
			err := json.Unmarshal(moveEvent.Message.Payload, &payload)
			if err != nil {
				logger.Error().Err(err).Send()
				continue
			}
			gameState, err := h.gameService.MakeMove(ctx, client.UserID, client.GameID, payload.Moves)
			if err != nil {
				logger.Error().Err(err).Send()
				continue
			}

			whitePlayerClient, err := h.findClientByID(gameState.WhitePlayerID)
			if err != nil {
				logger.Error().Err(err).Send()
				continue
			}

			blackPlayerClient, err := h.findClientByID(gameState.BlackPlayerID)
			if err != nil {
				logger.Error().Err(err).Send()
				continue
			}

			stateMsg, _ := json.Marshal(Message{
				Type:    MessageTypeGameState,
				Payload: json.RawMessage("Wait"),
			})

			whitePlayerClient.Send <- stateMsg
			blackPlayerClient.Send <- stateMsg
		}

	}
}

func (h *Hub) findClientByID(userID uuid.UUID) (*Client, error) {
	if _, ok := h.Clients[userID]; ok {
		return h.Clients[userID], nil
	}
	return nil, errors.New("user offline")
}

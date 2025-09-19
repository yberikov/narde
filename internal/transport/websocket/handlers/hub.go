package handlers

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
	"log"
	"narde/internal/transport/http/router"
	websocket2 "narde/internal/transport/websocket"
	"narde/pkg/jwt"
)

type HubHandlers struct {
	makeRouter router.MakeRouter
	parser     *jwt.Parser
	hub        *websocket2.Hub
}

func NewHubHander(parser *jwt.Parser, hub *websocket2.Hub) *HubHandlers {
	hubHandlers := &HubHandlers{parser: parser, hub: hub}
	makeRouter := func(router fiber.Router) {
		g := router.Group("/hub")
		g.Get("/ws", websocket.New(hubHandlers.wsHandler))
	}
	hubHandlers.makeRouter = makeRouter
	return hubHandlers
}

func (h *HubHandlers) Router() router.MakeRouter {
	return h.makeRouter
}

func (h *HubHandlers) wsHandler(c *websocket.Conn) {
	ctx := context.Background()
	logger := zerolog.Ctx(ctx)
	tokenStr := c.Query("token")
	if tokenStr == "" {
		log.Println("Closing connection: Missing auth token in query param")
		c.Close()
		return
	}

	sessionUser, err := h.parser.ParseToken(context.Background(), tokenStr)
	if err != nil {
		log.Printf("Closing connection: Invalid token for %s", c.RemoteAddr())
		c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInvalidFramePayloadData, "Invalid auth token"))
		c.Close()
		return
	}

	client := &websocket2.Client{
		Hub:    h.hub,
		Conn:   c,
		Send:   make(chan []byte, 256),
		UserID: sessionUser.ID,
	}

	logger.Info().Msgf("user %s added to hub", client.UserID.String())

	h.hub.Register <- client
	go client.WritePump()
	client.ReadPump()
}

package websocket

type MoveEvent struct {
	Client  *Client
	Message Message
}

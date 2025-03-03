package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"

	"github.com/Natnael-Alemayehu/chat/chat/app/sdk/errs"
	"github.com/Natnael-Alemayehu/chat/chat/foundation/logger"
)

type Chat struct {
	log *logger.Logger
}

func New(log *logger.Logger) *Chat {
	return &Chat{
		log: log,
	}
}

func (c *Chat) Handshake(ctx context.Context, conn *websocket.Conn) (User, error) {

	c.log.Info(ctx, "Chat", "Status", "Handshaking")
	defer c.log.Info(ctx, "Chat", "Status", "Finished handshaking")

	err := conn.WriteMessage(websocket.TextMessage, []byte("Hello"))
	if err != nil {
		return User{}, errs.Newf(errs.FailedPrecondition, "Failed to write Hello message: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	msg, err := c.readMessage(ctx, conn)
	if err != nil {
		return User{}, errs.Newf(errs.FailedPrecondition, "Failed to get : %v", err)
	}

	var usr User
	if err = json.Unmarshal(msg, &usr); err != nil {
		return User{}, errs.Newf(errs.FailedPrecondition, "Failed to unmarshal user: %v", err)
	}

	welcomeMessage := fmt.Sprintf("Welcome %s", usr.Name)
	if err = conn.WriteMessage(websocket.TextMessage, []byte(welcomeMessage)); err != nil {
		return User{}, errs.Newf(errs.FailedPrecondition, "Failed to write Welcome message: %v", err)
	}

	return usr, nil
}

func (c *Chat) readMessage(ctx context.Context, conn *websocket.Conn) ([]byte, error) {

	type response struct {
		msg []byte
		err error
	}

	ch := make(chan response, 1)

	go func() {
		c.log.Info(ctx, "Reading message")
		defer c.log.Info(ctx, "Finished reading message")
		_, msg, err := conn.ReadMessage()
		if err != nil {
			ch <- response{nil, err}
		}
		ch <- response{msg, nil}
	}()

	var resp response

	select {
	case <-ctx.Done():
		conn.Close()
		return nil, ctx.Err()
	case resp = <-ch:
		if resp.err != nil {
			return nil, errs.Newf(errs.FailedPrecondition, "Failed to read message")
		}
	}

	return resp.msg, nil
}

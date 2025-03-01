package chatapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/Natnael-Alemayehu/chat/chat/app/sdk/errs"
	"github.com/Natnael-Alemayehu/chat/chat/foundation/web"
)

type app struct {
	WS websocket.Upgrader
}

func newApp() *app {
	return &app{}
}

func (a *app) connect(ctx context.Context, r *http.Request) web.Encoder {

	c, err := a.WS.Upgrade(web.GetWriter(ctx), r, nil)
	if err != nil {
		return errs.Newf(errs.FailedPrecondition, "failed to connect to web socket: %v", err)
	}
	defer c.Close()

	_, err = a.handshake(ctx, c)
	if err != nil {
		return errs.Newf(errs.FailedPrecondition, "Failed handshake: %v", err)
	}

	return status{
		Status: "OK",
	}
}

func (a *app) handshake(ctx context.Context, c *websocket.Conn) (user, error) {

	err := c.WriteMessage(websocket.TextMessage, []byte("Hello"))
	if err != nil {
		return user{}, errs.Newf(errs.FailedPrecondition, "Failed to write Hello message: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	msg, err := readMessage(ctx, c)
	if err != nil {
		return user{}, errs.Newf(errs.FailedPrecondition, "Failed to get : %v", err)
	}

	var usr user
	if err = json.Unmarshal(msg, usr); err != nil {
		return user{}, errs.Newf(errs.FailedPrecondition, "Failed to unmarshal user: %v", err)
	}

	welcomeMessage := fmt.Sprintf("Welcome %s", usr.name)
	if c.WriteMessage(websocket.TextMessage, []byte(welcomeMessage)); err != nil {
		return user{}, errs.Newf(errs.FailedPrecondition, "Failed to write Welcome message: %v", err)
	}

	return usr, nil
}

func readMessage(ctx context.Context, c *websocket.Conn) ([]byte, error) {

	type response struct {
		msg []byte
		err error
	}

	ch := make(chan response, 1)

	go func() {
		_, msg, err := c.ReadMessage()
		if err != nil {
			ch <- response{nil, err}
		}
		ch <- response{msg, nil}
	}()

	var resp response

	select {
	case <-ctx.Done():
		c.Close()
		return nil, ctx.Err()
	case resp = <-ch:
		if resp.err != nil {
			return nil, errs.Newf(errs.FailedPrecondition, "Failed to read message")
		}
	}

	return resp.msg, nil
}

package chatapp

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/Natnael-Alemayehu/chat/chat/app/sdk/chat"
	"github.com/Natnael-Alemayehu/chat/chat/app/sdk/errs"
	"github.com/Natnael-Alemayehu/chat/chat/foundation/logger"
	"github.com/Natnael-Alemayehu/chat/chat/foundation/web"
)

type app struct {
	log  *logger.Logger
	WS   websocket.Upgrader
	chat *chat.Chat
}

func newApp(log *logger.Logger) *app {
	return &app{
		log:  log,
		chat: chat.New(log),
	}
}

func (a *app) connect(ctx context.Context, r *http.Request) web.Encoder {

	c, err := a.WS.Upgrade(web.GetWriter(ctx), r, nil)
	if err != nil {
		return errs.Newf(errs.FailedPrecondition, "failed to connect to web socket: %v", err)
	}
	defer c.Close()

	_, err = a.chat.Handshake(ctx, c)
	if err != nil {
		return errs.Newf(errs.FailedPrecondition, "Failed handshake: %v", err)
	}

	return web.NewNoResponse()
}

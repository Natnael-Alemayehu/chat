package chatapp

import (
	"net/http"

	"github.com/Natnael-Alemayehu/chat/chat/foundation/logger"
	"github.com/Natnael-Alemayehu/chat/chat/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App, log *logger.Logger) {
	const version = "v1"

	api := newApp(log)

	app.HandlerFunc(http.MethodGet, version, "/connect", api.connect)

}

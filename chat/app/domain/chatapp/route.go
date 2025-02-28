package chatapp

import (
	"net/http"

	"github.com/Natnael-Alemayehu/chat/chat/foundation/web"
)

// Routes adds specific routes for this group.
func Routes(app *web.App) {
	const version = "v1"

	api := newApp()

	app.HandlerFunc(http.MethodGet, version, "/status", api.test)

}

package main

import (
	"net/http"

	"github.com/jcroyoaun/snackmath/assets"

	"github.com/alexedwards/flow"
)

func (app *application) routes() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(app.notFound)

	mux.Use(app.recoverPanic)
	mux.Use(app.securityHeaders)

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.Handle("/static/...", fileServer, "GET")

	mux.HandleFunc("/sw.js", app.serviceWorker, "GET")
	mux.HandleFunc("/", app.home, "GET")

	return mux
}

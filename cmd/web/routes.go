package main

import (
	"net/http"

	"example/ui"

	"github.com/justinas/alice"
)

// Update the signature routes method so that it returns a http.Handler instead of *http.ServerMux
func (app *application) routes() http.Handler {

	// Use the http.NewServeMux to initialize a new servermux
	mux := http.NewServeMux()

	// Create a file server which serves files out of the directory
	//fileServer := http.FileServer(http.Dir("../../ui/static/"))
	// Register the file server as the handler for all url paths that start with "/static/"
	// http://localhost:4000/static
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	// Register the home function as the handler for the "/" url pattern
	// {$} restrict this route to exact matches on "/" only
	// Create a new middleware chain containing the middleware specific to our app
	//dynamic := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	// Use the nosurf middleware on all our dynamic routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// https://localhost:4000/
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	// https://localhost:4000/snippet/view/1
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))

	// https://localhost:4000/snippet/create
	//mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	//mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
	// https://localhost:4000/user/logout
	//mux.Handle("POST /user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// https://localhost:4000/user/signup
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	//
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))

	// https://localhost:4000/user/login
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	//
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// https://localhost:4000/about
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))

	protected := dynamic.Append(app.requireAuthentication)
	// https://localhost:4000/snippet/create
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))

	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// https://localhost:4000/snippet/update/1
	mux.Handle("GET /snippet/update/{id}", protected.ThenFunc(app.snippetUpdate))
	mux.Handle("POST /snippet/update/{id}", protected.ThenFunc(app.snippetUpdatePost))

	// https://localhost:4000/snippet/delete/1
	mux.Handle("GET /snippet/delete/{id}", protected.ThenFunc(app.snippetDelete))
	mux.Handle("POST /snippet/delete/{id}", protected.ThenFunc(app.snippetDeletePost))

	// https://localhost:4000/account/view
	mux.Handle("GET /account/view", protected.ThenFunc(app.accountView))
	// https://localhost:4000/account/password/update
	mux.Handle("GET /account/password/update", protected.ThenFunc(app.accountPassUpdate))
	mux.Handle("POST /account/password/update", protected.ThenFunc(app.accountPassUpdatePost))

	standar := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	// Pass the servermux as the next parameter to the commonHeaders middleware
	// Wrap the existing chain with the logRequest middleware
	// Wrap the existing chain with the recoverPanic middleware
	//return app.recoverPanic(app.logRequest(commonHeaders(mux)))
	return standar.Then(mux)
}

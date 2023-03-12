package main

import (
	"net/http"

	"mathrix.ceg.com/ui"
	
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.FS(ui.Files))

	router.Handler(http.MethodGet, "/static/*filepath",fileServer)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.validateAuthority)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/event", dynamic.ThenFunc(app.event))
	router.Handler(http.MethodGet, "/event/view/:id", dynamic.ThenFunc(app.eventView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/event/add/:id", protected.ThenFunc(app.eventAddPost))
	
	authorized:=protected.Append(app.requireAuthority)

	router.Handler(http.MethodGet, "/user", authorized.ThenFunc(app.user))
	router.Handler(http.MethodGet, "/user/create", authorized.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/register", authorized.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/update/:id", authorized.ThenFunc(app.userUpdate))
	router.Handler(http.MethodGet, "/user/delete/:id", authorized.ThenFunc(app.userDelete))
	router.Handler(http.MethodGet, "/event/create", authorized.ThenFunc(app.eventCreate))
	router.Handler(http.MethodPost, "/event/create", authorized.ThenFunc(app.eventCreatePost))
	router.Handler(http.MethodGet, "/event/update/:id", authorized.ThenFunc(app.eventUpdate))
	router.Handler(http.MethodPost, "/event/update/:id", authorized.ThenFunc(app.eventUpdatePost))
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}

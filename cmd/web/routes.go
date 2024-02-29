package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/news", dynamicMiddleware.ThenFunc(app.showNews))
	mux.Get("/news/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNews))
	mux.Post("/news/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createNews))
	mux.Get("/news/creationPage", dynamicMiddleware.ThenFunc(app.creationPage))
	mux.Get("/news/contacts", dynamicMiddleware.ThenFunc(app.contacts))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/role", dynamicMiddleware.ThenFunc(app.Role))
	mux.Get("/user/admin", dynamicMiddleware.ThenFunc(app.Admin))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))
	mux.Get("/news/students", dynamicMiddleware.ThenFunc(app.filterCategory("Students")))
	mux.Get("/news/applicants", dynamicMiddleware.ThenFunc(app.filterCategory("Applicants")))
	mux.Get("/news/researches", dynamicMiddleware.ThenFunc(app.filterCategory("Researches")))
	mux.Get("/news/staff", dynamicMiddleware.ThenFunc(app.filterCategory("Staff")))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}

package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))

	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()

	secure := alice.New(app.checkToken)

	router.HandlerFunc(http.MethodPost, "/v1/signup", app.signUp)
    router.HandlerFunc(http.MethodPost, "/v1/signin", app.signIn)

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)


	router.HandlerFunc(http.MethodGet, "/v1/thread/:id", app.getOneThread)
	router.HandlerFunc(http.MethodGet, "/v1/threads", app.getAllThreads)
	router.HandlerFunc(http.MethodGet, "/v1/threads/:category_id", app.getAllThreadsByCategory)
	router.HandlerFunc(http.MethodGet, "/v1/categories", app.getAllCategories)
	router.HandlerFunc(http.MethodGet, "/v1/replies/:thread_id", app.getReplies)

	router.POST("/v1/admin/editthread", app.wrap(secure.ThenFunc(app.editThread)))
	//router.HandlerFunc(http.MethodPost, "/v1/admin/editthread", app.editThread)

	router.GET("/v1/admin/deletethread/:id", app.wrap(secure.ThenFunc(app.deleteThread)))
	//router.HandlerFunc(http.MethodGet, "/v1/admin/deletethread/:id", app.deleteThread)

	router.HandlerFunc(http.MethodPost, "/v1/newreply/:thread_id", app.newReply)
	router.HandlerFunc(http.MethodGet, "/v1/deletereply/:id", app.deleteReply)
	router.HandlerFunc(http.MethodPost, "/v1/newthread", app.newThread)
	router.HandlerFunc(http.MethodPost, "/v1/editthread/", app.editThread)
	router.HandlerFunc(http.MethodPut, "/v1/togglesolved/:id", app.toggleSolved)
	router.HandlerFunc(http.MethodPut, "/v1/toggleanswer/:reply_id", app.toggleAnswer)
	router.HandlerFunc(http.MethodGet, "/v1/deletethread/:id", app.deleteThread)
	router.HandlerFunc(http.MethodGet, "/v1/yourthreads/:author_id", app.yourThreads)
	
	router.HandlerFunc(http.MethodPost, "/v1/star/:user_id/:thread_id", app.starThread)
    router.HandlerFunc(http.MethodDelete, "/v1/unstar/:user_id/:thread_id", app.unstarThread)
    router.HandlerFunc(http.MethodGet, "/v1/starred/:user_id", app.getStarredThreads)
	
	return app.enableCORS(router)

}

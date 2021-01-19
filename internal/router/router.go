package router

import (
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/internal/controllers"
	"github.com/semihsemih/save-my-notes/internal/middleware"
	"github.com/urfave/negroni"
)

func Init(controller *controllers.Controller) *mux.Router {
	/* Main router */
	router := mux.NewRouter()
	/* API router */
	apiRouter := mux.NewRouter()
	/* Auth router */
	authRouter := mux.NewRouter()

	/* User Action Routes */
	apiRouter.HandleFunc("/api/user/{id:[0-9]+}", controller.GetUser()).Methods("GET")

	/* List Action Routes */
	apiRouter.HandleFunc("/api/list", controller.InsertList()).Methods("POST")
	apiRouter.HandleFunc("/api/list/{id:[0-9]+}", controller.GetList()).Methods("GET")
	apiRouter.HandleFunc("/api/list/{id:[0-9]+}", controller.UpdateList()).Methods("PUT")
	apiRouter.HandleFunc("/api/list/{id:[0-9]+}", controller.DeleteList()).Methods("DELETE")

	/* Note Action Routes */
	apiRouter.HandleFunc("/api/note", controller.InsertNote()).Methods("POST")
	apiRouter.HandleFunc("/api/note/{id:[0-9]+}", controller.GetNote()).Methods("GET")
	apiRouter.HandleFunc("/api/note/{id:[0-9]+}", controller.UpdateNote()).Methods("PUT")
	apiRouter.HandleFunc("/api/note/{id:[0-9]+}", controller.DeleteNote()).Methods("DELETE")

	/* Authentication Action Routes */
	authRouter.HandleFunc("/auth/signup", controller.Signup()).Methods("POST")
	authRouter.HandleFunc("/auth/login", controller.Login()).Methods("POST")
	authRouter.HandleFunc("/auth/activation/{token}", controller.AccountActivation()).Methods("GET")

	/* API routes middlewares */
	router.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(middleware.CORS),
		negroni.HandlerFunc(middleware.TokenVerifyMiddleware),
		negroni.HandlerFunc(middleware.GzipMiddleware),
		negroni.Wrap(apiRouter),
	))

	/* Auth routes middlewares */
	router.PathPrefix("/auth").Handler(negroni.New(
		negroni.HandlerFunc(middleware.CORS),
		negroni.HandlerFunc(middleware.GzipMiddleware),
		negroni.Wrap(authRouter),
	))

	return router
}

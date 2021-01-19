package router

import (
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/internal/controllers"
	"github.com/semihsemih/save-my-notes/internal/middleware"
	"github.com/urfave/negroni"
)

func Init(controller *controllers.Controller) *mux.Router {
	router := mux.NewRouter()
	apiBaseRouter := mux.NewRouter()

	/* Subrouters */
	apiRouter := apiBaseRouter.PathPrefix("/api").Subrouter()
	authRouter := router.PathPrefix("/auth").Subrouter()

	/* Router middlwares */
	router.PathPrefix("/").Handler(negroni.New(
		negroni.HandlerFunc(middleware.GzipMiddleware),
		negroni.HandlerFunc(middleware.CORS),
	))

	router.PathPrefix("/api").Handler(negroni.New(
		negroni.HandlerFunc(middleware.TokenVerifyMiddleware),
		negroni.Wrap(router),
	))

	/* Authentication Action Routes */
	authRouter.HandleFunc("/signup", controller.Signup()).Methods("POST")
	authRouter.HandleFunc("/login", controller.Login()).Methods("POST")
	authRouter.HandleFunc("/activation/{token}", controller.AccountActivation()).Methods("GET")

	/* User Action Routes */
	apiRouter.HandleFunc("/user/{id:[0-9]+}", controller.GetUser()).Methods("GET")

	/* List Action Routes */
	apiRouter.HandleFunc("/list", controller.InsertList()).Methods("POST")
	apiRouter.HandleFunc("/list/{id:[0-9]+}", controller.GetList()).Methods("GET")
	apiRouter.HandleFunc("/list/{id:[0-9]+}", controller.UpdateList()).Methods("PUT")
	apiRouter.HandleFunc("/list/{id:[0-9]+}", controller.DeleteList()).Methods("DELETE")

	/* Note Action Routes */
	apiRouter.HandleFunc("/note", controller.InsertNote()).Methods("POST")
	apiRouter.HandleFunc("/note/{id:[0-9]+}", controller.GetNote()).Methods("GET")
	apiRouter.HandleFunc("/note/{id:[0-9]+}", controller.UpdateNote()).Methods("PUT")
	apiRouter.HandleFunc("/note/{id:[0-9]+}", controller.DeleteNote()).Methods("DELETE")

	return router
}
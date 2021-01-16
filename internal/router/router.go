package router

import (
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/internal/controllers"
)

func Init(controller *controllers.Controller) *mux.Router {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()
	authRouter := router.PathPrefix("/auth").Subrouter()

	/* Authentication Action Routes */
	authRouter.HandleFunc("/signup", controller.Signup()).Methods("POST")
	authRouter.HandleFunc("/login", controller.Login()).Methods("POST")
	authRouter.HandleFunc("/activation/{token}", controller.AccountActivation()).Methods("GET")

	/* User Action Routes */
	apiRouter.HandleFunc("/user/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.GetUser())).Methods("GET")

	/* List Action Routes */
	apiRouter.HandleFunc("/list", controller.TokenVerifyMiddleware(controller.InsertList())).Methods("POST")
	apiRouter.HandleFunc("/list/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.GetList())).Methods("GET")
	apiRouter.HandleFunc("/list/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.UpdateList())).Methods("PUT")
	apiRouter.HandleFunc("/list/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.DeleteList())).Methods("DELETE")

	/* Note Action Routes */
	apiRouter.HandleFunc("/note", controller.TokenVerifyMiddleware(controller.InsertNote())).Methods("POST")
	apiRouter.HandleFunc("/note/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.GetNote())).Methods("GET")
	apiRouter.HandleFunc("/note/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.UpdateNote())).Methods("PUT")
	apiRouter.HandleFunc("/note/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.DeleteNote())).Methods("DELETE")

	return router
}
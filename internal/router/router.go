package router

import (
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/internal/controllers"
)

func Init(controller *controllers.Controller) *mux.Router {
	router := mux.NewRouter()
	/* User Action Routes */
	router.HandleFunc("/api/user/signup", controller.Signup()).Methods("POST")
	router.HandleFunc("/api/user/login", controller.Login()).Methods("POST")
	router.HandleFunc("/api/user/activation/{token}", controller.AccountActivation()).Methods("GET")
	router.HandleFunc("/api/user/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.GetUser())).Methods("GET")


	/* List Action Routes */
	router.HandleFunc("/api/list", controller.TokenVerifyMiddleware(controller.InsertList())).Methods("POST")
	router.HandleFunc("/api/list/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.GetList())).Methods("GET")
	router.HandleFunc("/api/list/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.UpdateList())).Methods("PUT")
	router.HandleFunc("/api/list/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.DeleteList())).Methods("DELETE")

	/* Note Action Routes */
	router.HandleFunc("/api/note", controller.TokenVerifyMiddleware(controller.InsertNote())).Methods("POST")
	router.HandleFunc("/api/note/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.GetNote())).Methods("GET")
	router.HandleFunc("/api/note/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.UpdateNote())).Methods("PUT")
	router.HandleFunc("/api/note/{id:[0-9]+}", controller.TokenVerifyMiddleware(controller.DeleteNote())).Methods("DELETE")

	return router
}
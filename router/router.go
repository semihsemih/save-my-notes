package router

import (
	"github.com/gorilla/mux"
	"github.com/semihsemih/save-my-notes/controllers"
	"github.com/semihsemih/save-my-notes/driver"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() *mux.Router {
	db = driver.ConnectDB()
	controller := controllers.Controller{}

	router := mux.NewRouter()
	/* User Action Routes */
	router.HandleFunc("/api/user/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/api/user/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/api/user/activation/{token}", controller.AccountActivation(db)).Methods("GET")
	router.HandleFunc("/api/user/{id}", controller.TokenVerifyMiddleware(controller.GetUser(db))).Methods("GET")


	/* List Action Routes */
	router.HandleFunc("/api/list", controller.TokenVerifyMiddleware(controller.InsertList(db))).Methods("POST")
	router.HandleFunc("/api/list/{id}", controller.TokenVerifyMiddleware(controller.GetList(db))).Methods("GET")
	router.HandleFunc("/api/list/{id}", controller.TokenVerifyMiddleware(controller.UpdateList(db))).Methods("PUT")
	router.HandleFunc("/api/list/{id}", controller.TokenVerifyMiddleware(controller.DeleteList(db))).Methods("DELETE")

	/* Note Action Routes */
	router.HandleFunc("/api/note", controller.TokenVerifyMiddleware(controller.InsertNote(db))).Methods("POST")
	router.HandleFunc("/api/note/{id}", controller.TokenVerifyMiddleware(controller.GetNote(db))).Methods("GET")
	router.HandleFunc("/api/note/{id}", controller.TokenVerifyMiddleware(controller.UpdateNote(db))).Methods("PUT")
	router.HandleFunc("/api/note/{id}", controller.TokenVerifyMiddleware(controller.DeleteNote(db))).Methods("DELETE")

	return router
}
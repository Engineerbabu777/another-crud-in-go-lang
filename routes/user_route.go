package routes

import (
	"example.com/controllers"
	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router)  {
    //All routes related to users comes here
	router.HandleFunc("/user", controllers.CreateUser()).Methods("POST") //add this

	router.HandleFunc("/user/{userId}", controllers.GetAUser()).Methods("GET") //add this

	router.HandleFunc("/user/{userId}", controllers.EditAUser()).Methods("PUT") //add this

	router.HandleFunc("/user/{userId}", controllers.DeleteAUser()).Methods("DELETE") //add this

	router.HandleFunc("/users", controllers.GetAllUser()).Methods("GET") //add this

}
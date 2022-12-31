package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-telegram-bot-api/controllers"
	"github.com/gorilla/mux"
)

func InitialiseApi() {
	r := mux.NewRouter()
	homeController := controllers.NewHomeController()
	r.HandleFunc("/", homeController.Home).Methods("GET")
	port := os.Getenv("PORT")
	fmt.Println("running on port " + port)

	http.ListenAndServe(":"+port, r)
}

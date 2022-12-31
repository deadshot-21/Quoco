package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HomeController struct {
}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (hc HomeController) Home(w http.ResponseWriter, r *http.Request) {
	// var credentials Credentials
	// err := json.NewDecoder(r.Body).Decode(&credentials)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

	vars := make(map[string]string)
	json.NewDecoder(r.Body).Decode(&vars)
	// username := vars["username"]
	// password := vars["password"]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", `{
		"success": "true",
		"message": "Welcome to Quoco!",
		"data": {},
	}`)

}

package main

import (
	"Login-with-Outh/controller"
	"net/http"
)

func main() {
	http.HandleFunc("/google/login", controller.GoogleLogin)
	http.HandleFunc("/google/callback", controller.GoogleCallback)

	http.ListenAndServe(":8080", nil)
}

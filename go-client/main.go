package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", loadAuthPage)
	mux.HandleFunc("/login-google", loginGoogle)

	server := &http.Server{
		Handler: mux,
		Addr: "0.0.0.0:8080",
	}

	fmt.Printf("start server ... ")
	panic(server.ListenAndServe())
}

func loadAuthPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/login.html")
}

func loginGoogle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "reached google login page")
}


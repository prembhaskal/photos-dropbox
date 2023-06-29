package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func main() {
	conf := &oauth2.Config{
		ClientID:     "", // TODO - get from env. / json file
		ClientSecret: "", // TODO - get from env. / json file
		Scopes: []string{
			"https://www.googleapis.com/auth/drive.metadata.readonly",
			"https://www.googleapis.com/auth/photoslibrary.readonly"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://oauth2.googleapis.com/token", // got this from json downloaded from console
		},
		RedirectURL: "http://localhost:19080/callback",
	}

	mux := &http.ServeMux{}
	mux.HandleFunc("/", loadAuthPage)
	mux.HandleFunc("/login-google", loginGoogle(conf))
	mux.HandleFunc("/callback", callbackHandler(conf))

	server := &http.Server{
		Handler: mux,
		Addr:    "0.0.0.0:8080",
	}

	fmt.Printf("start server ... ")
	panic(server.ListenAndServe())
}

func loadAuthPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/login.html")
}

func callbackHandler(conf *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got request: %+v in callback\n", r)

		queryParams := r.URL.Query()

		code := queryParams.Get("code")
		log.Printf("in callback, got code: %s", code)
		if len(code) == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: return auth code is tmpy")
			return
		}

		tok, err := conf.Exchange(context.Background(), code)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error in getting token: %v", err)
			return
		}

		log.Printf("token obtained: %+v", tok)

	}
}

func loginGoogle(conf *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "reached google login page")

		conf := &oauth2.Config{
			ClientID:     "1082052306771-igus2n4c7tqof0umof9s5i7o9ant9p1p.apps.googleusercontent.com",
			ClientSecret: "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",
			Scopes: []string{
				"https://www.googleapis.com/auth/drive.metadata.readonly",
				"https://www.googleapis.com/auth/photoslibrary.readonly"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
				TokenURL: "https://accounts.google.com/o/oauth2/v2/token",
			},
			RedirectURL: "http://localhost:19080/callback",
		}

		url := conf.AuthCodeURL("test123", oauth2.AccessTypeOffline)

		log.Printf("redirecting to url: %v for getting code\n", url)

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

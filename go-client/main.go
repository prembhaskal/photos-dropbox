package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"io"

	"github.com/jeremywohl/flatten"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
)

type PhotosConfig struct {
	TokenFile string `envconfig:"TOKEN_FILE" default:"/var/tmp/client_info.json"`
}

func main() {
	clientID, clientSecret, err := getClientSecret()
	if err != nil {
		panic(fmt.Sprintf("unable to read clientID or secret, check env variable TOKEN_FILE, err: %v", err))
	}

	tokenChan := make(chan string, 0)
	go getToken(clientID, clientSecret, tokenChan)

	var token string
	select {
	case token = <-tokenChan:
		break
	case <-time.After(1 * time.Minute):
		log.Printf("timed out obtaining token, exiting...")
		panic("timed out token")
	}

	log.Printf("processing with token: %s", token)
	err = listPhotos("", token)
	if err != nil {
		log.Printf("Error in list pics: %v", err)
	}
}

func listPhotos(nextpagetoken, accesstoken string) error {
	listPicsURL := "https://photoslibrary.googleapis.com/v1/mediaItems?pageSize=10&pageToken=" + nextpagetoken
	req, err := http.NewRequest("GET", listPicsURL, nil)
	if err != nil {
		return err
	}

	header := http.Header{}
	header.Add("Authorization", fmt.Sprintf("Bearer %s", accesstoken))

	req.Header = header

	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Printf("list photos output: %+v", resp)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Printf("list pics resp body: %s", string(data))

	return nil
}

func getClientSecret() (string, string, error) {
	conf := &PhotosConfig{}
	err := envconfig.Process("", conf)
	if err != nil {
		return "", "", err
	}

	// read the file
	data, err := os.ReadFile(conf.TokenFile)
	if err != nil {
		return "", "", fmt.Errorf("error reading file: %s, err: %w", conf.TokenFile, err)
	}

	var jsonval map[string]any

	err = json.Unmarshal(data, &jsonval)
	if err != nil {
		return "", "", err
	}

	nested, err := flatten.Flatten(jsonval, "", flatten.DotStyle)
	if err != nil {
		return "", "", err
	}

	clientID := nested["web.client_id"].(string)
	clientSecret := nested["web.client_secret"].(string)
	if len(clientID) == 0 {
		return "", "", fmt.Errorf("empty clientID")
	}
	if len(clientSecret) == 0 {
		return "", "", fmt.Errorf("empty client secret")
	}

	return clientID, clientSecret, nil

}

func getToken(clientID, clientSecret string, tokenChan chan<- string) {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
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
	mux.HandleFunc("/callback", callbackHandler(conf, tokenChan))

	server := &http.Server{
		Handler: mux,
		Addr:    "0.0.0.0:8080",
	}

	log.Println("starting server ...")
	log.Println("open localhost:8080 and follow steps to get token")

	panic(server.ListenAndServe())
}

func loadAuthPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/login.html")
}

func callbackHandler(conf *oauth2.Config, tokenChan chan<- string) http.HandlerFunc {
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
		tokenChan <- tok.AccessToken
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

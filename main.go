package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	PORT         = "6161"
	ClientID     = os.Getenv("ClientID")     //NOTE:  Client ID of oauth app in github
	ClientSecret = os.Getenv("ClientSecret") //NOTE: Client Secret of oauth app in github
)

func main() {

	r := http.NewServeMux()

	r.Handle("/", http.FileServer(http.Dir("./static")))

	r.HandleFunc("/gh/login", ghLogin)

	r.HandleFunc("/gh/callback", ghCallback)

	r.HandleFunc("/login/successfull", successHandler)

	log.Printf("Server starting at %s", PORT)

	log.Fatal(http.ListenAndServe(":"+PORT, r))

}

func successHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Logged in"))

}

func ghLogin(w http.ResponseWriter, r *http.Request) {

	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		ClientID,
		"http://localhost:6161/gh/callback",
	)

	http.Redirect(w, r, redirectURL, 301)

}

func ghCallback(w http.ResponseWriter, r *http.Request) {

	code := r.URL.Query().Get("code")

	token, err := GetAccessToken(code)

	if err != nil {

		log.Println("Error obtaining token from gh: ", err)

		return
	}

	data, err := GetUserData(token)

	//NOTE: User data to be used in other operations

	log.Println(data)

	if err != nil {

		log.Println("Error getting user data: ", err)

		return
	}

	http.Redirect(w, r, "/login/successfull", http.StatusTemporaryRedirect)
}

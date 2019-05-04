package main

import (
	"github.com/jacobtie/go-chat-room/backend/internal/pkg/auth"

	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func renderTemplate(w http.ResponseWriter, filename string, data map[string]interface{}) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.ParseFiles(wd+"/build/base.html", wd+"/build/"+filename)
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, "base", data)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET info")
	renderTemplate(w, "index.html", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("POST login")
	r.ParseForm()
	if password := r.FormValue("password"); password == "greenlantern" {
		log.Println("Login successful")
		validToken, err := auth.GenerateJWT(r.FormValue("username"))
		if err != nil {
			http.Error(w, "Error generating token", 500)
		}
		http.Redirect(w, r, "/chat/"+validToken, 301)
	} else {
		log.Println("Login unsuccessful, " + r.FormValue("password"))
		http.Redirect(w, r, "/", 301)
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET chat")
	renderTemplate(w, "chat.html", nil)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler).Methods("GET")
	r.HandleFunc("/chat/{token}", auth.MustAuth(chatHandler)).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting server on PORT " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

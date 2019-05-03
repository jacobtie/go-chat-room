package main

import (
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler).Methods("GET")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting server on PORT " + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

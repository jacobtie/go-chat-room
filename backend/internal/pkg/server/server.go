package server

import (
	"github.com/jacobtie/go-chat-room/backend/internal/pkg/auth"
	"github.com/jacobtie/go-chat-room/backend/internal/pkg/ws"

	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

var hub *ws.Hub

func setUpRoutes(r *mux.Router) {
	r.HandleFunc("/", mainHandler).Methods("GET")
	r.HandleFunc("/chat", auth.MustAuth(chatHandler)).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/ws", auth.MustAuth(wsHandler))
}

func renderTemplate(w http.ResponseWriter, filename string, data map[string]interface{}) {
	wd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.ParseFiles(wd+"/base.html", wd+"/"+filename)
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
	if password := r.FormValue("password"); password == os.Getenv("GO_CHAT_PASS") {
		log.Println("Login successful")
		validToken, err := auth.GenerateJWT()
		if err != nil {
			http.Error(w, "Error generating token", 500)
		}
		c := &http.Cookie{
			Name:     "jwt",
			Value:    validToken,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(time.Minute * 30),
		}
		http.SetCookie(w, c)
		userCookie := &http.Cookie{
			Name:    "username",
			Value:   r.FormValue("username"),
			Secure:  true,
			Expires: time.Now().Add(time.Minute * 30),
		}
		http.SetCookie(w, userCookie)
		http.Redirect(w, r, "/chat", 301)
	} else {
		log.Println("Login unsuccessful, " + r.FormValue("password"))
		http.Redirect(w, r, "/", 301)
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET chat")
	renderTemplate(w, "chat.html", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws.ServeWs(hub, w, r)
}

// Run starts the server
func Run() {
	hub = ws.NewHub()
	go hub.Run()

	r := mux.NewRouter()
	setUpRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Starting server on PORT " + port)
	wd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServeTLS(":"+port, wd+"/localhost.pem", wd+"/localhost-key.pem", r))
	// log.Fatal(http.ListenAndServe(":"+port, r))
}

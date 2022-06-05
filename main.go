package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	auth "plc-backend/Auth"
	"plc-backend/Controllers"
	"plc-backend/File"
	"plc-backend/Helmet"
	s "plc-backend/Settings"
	ws "plc-backend/Websocket"

	"github.com/gorilla/mux"
	env "github.com/joho/godotenv"
	"github.com/rs/cors"
)

var settings = s.NewSettingsHandler()

var store = new(auth.Sessions).Init()

func init() {
	// Load environment variables
	err := env.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read the settings from storage
	if File.Exists(os.Getenv("SETTINGS")) {
		strSettings, _ := File.Read(os.Getenv("SETTINGS"))
		json.Unmarshal([]byte(strSettings), &settings.Settings)
	}
}

func main() {
	r := mux.NewRouter()

	r.Use(Helmet.New())

	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/usr/frontend/static"))))

	// Websocket handler
	r.HandleFunc("/ws", ws.Handler)

	// Handle unprotected routes
	r.HandleFunc("/users/signup", Controllers.UserRegister).Methods("POST")
	r.HandleFunc("/users/login", Controllers.UserLogin(store)).Methods("POST")

	// Protected routes
	// Create a subrouter
	protected := r.PathPrefix("/auth/").Subrouter()
	protected.Use(auth.CheckJwtValidity)

	// Handle protected routes

	// User management
	//protected.HandleFunc("/users/update", Controllers.UpdateUser).Methods("POST")
	protected.HandleFunc("/users/logout", Controllers.UserLogout(store)).Methods("POST")

	// Settings file management
	protected.HandleFunc("/settings", settings.Get).Methods("GET")
	protected.HandleFunc("/settings", settings.Set).Methods("POST")

	// Add CORS headers
	handler := cors.Default().Handler(r)

	// Serve
	log.Println("Starting server on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}

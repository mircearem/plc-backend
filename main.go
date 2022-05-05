package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"plc-backend/File"
	"plc-backend/Helmet"
	"plc-backend/Routes"
	s "plc-backend/Utils"
	ws "plc-backend/Websocket"

	"github.com/gorilla/mux"
	env "github.com/joho/godotenv"
	"github.com/rs/cors"
)

var settings = s.NewSettingsHandler()

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

func setup() {
	r := mux.NewRouter()

	r.Use(Helmet.New())

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/usr/frontend/static"))))

	// Websocket handler
	r.HandleFunc("/ws", ws.Write)

	// Route handlers
	r.HandleFunc("/users/signup", Routes.Register).Methods("POST")
	r.HandleFunc("/users/update", Routes.UpdateUser).Methods("POST")
	r.HandleFunc("/users/login", Routes.Login).Methods("POST")
	r.HandleFunc("/settings", settings.Get).Methods("GET")
	r.HandleFunc("/settings", settings.Set).Methods("POST")

	// Add CORS headers
	handler := cors.Default().Handler(r)

	// Serve
	log.Println("Starting server on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", handler))
}

func main() {
	setup()
}

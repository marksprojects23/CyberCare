package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

// Settings holds whitelist and blacklist settings for a user-
type Settings struct {
	Whitelist []string `json:"whitelist"`
	Blacklist []string `json:"blacklist"`
}

var db *sql.DB

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust this to your needs
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Check if it's a preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the request to the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	dbHost := "database-1.cl0i0y628wpg.us-east-1.rds.amazonaws.com"
	dbPort := 5432
	dbUser := "postgres"
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := "postgres"

	// Connect to the PostgreSQL database
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", dbHost, dbPort, dbUser, dbPassword, dbName)
	var err error
	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}

	r := mux.NewRouter()
	r.Use(enableCORS)
	r.HandleFunc("/settings", CreateSettings).Methods("POST")
	r.HandleFunc("/settings", GetSettings).Methods("GET")
	r.HandleFunc("/settings", UpdateSettings).Methods("PUT")

	http.ListenAndServe(":8080", r)
}

// CreateSettings handles POST requests to create new settings
func CreateSettings(w http.ResponseWriter, r *http.Request) {
	var newSettings Settings
	if err := json.NewDecoder(r.Body).Decode(&newSettings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert settings into the database
	_, err := db.Exec("INSERT INTO setting (whitelist, blacklist) VALUES ($1, $2)", pq.Array(newSettings.Whitelist), pq.Array(newSettings.Blacklist))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newSettings)
}

// GetSettings handles GET requests to retrieve all settings
func GetSettings(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT whitelist, blacklist FROM setting")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var settingsList []Settings
	for rows.Next() {
		var s Settings
		if err := rows.Scan(pq.Array(&s.Whitelist), pq.Array(&s.Blacklist)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		settingsList = append(settingsList, s)
	}

	json.NewEncoder(w).Encode(settingsList)
}

func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var settings Settings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Clear the existing settings
	_, err := db.Exec("DELETE FROM setting")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the new settings
	_, err = db.Exec("INSERT INTO setting (whitelist, blacklist) VALUES ($1, $2)",
		pq.Array(settings.Whitelist), pq.Array(settings.Blacklist))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Or http.StatusOK if you prefer
}

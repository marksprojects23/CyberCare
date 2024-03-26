package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
	_ "golang.org/x/crypto/ssh"
)

// Define a simple struct to hold settings for a user
type Settings struct {
    UserID    string   `json:"userId"`
    Whitelist []string `json:"whitelist"`
    Blacklist []string `json:"blacklist"`
}

// var (
// 	ctx context.Context
// 	db sql.DB
// )

// Temporary storage for settings
var settingsMap = make(map[string]Settings)

func main() {
	r := mux.NewRouter()
    r.HandleFunc("/settings", CreateSettings).Methods("POST")
    r.HandleFunc("/settings/{userId}", GetSettings).Methods("GET")

    http.ListenAndServe(":8080", r)

	// PostgreSQL server credentials
	// dbHost := "localhost"
	// dbPort := 5432
	// dbUser := "postgres"
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := "postgres"

// Connect to PostgreSQL database
	// connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", dbHost, dbPort, dbUser, dbPassword, dbName)
	// db, err := sql.Open("postgres", connString)
	// if err != nil {log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	// }
	// rows, err := db.Query ("SELECT * FROM users;")
	// log.Output(1, fmt.Sprintf("%v", rows))
	// if err != nil {
	// 	log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	// }
	// defer db.Close()

	// for rows.Next() {
	// 	var (
	// 		id  int
	// 		username string
	// 		passwordHash string
	// 		email string
	// 		createdAt time.Time
	// 		updatedAt time.Time

	// 	)
	// 	err := rows.Scan(&id, &username, &passwordHash, &email, &createdAt, &updatedAt)
	// 	if err != nil {log.Fatalf("Error scanning row: ", err)
	// 	}

	// 	fmt.Printf("ID: %v, Password: %v, Email: %s, CreatedAt: %v, UpdatedAt: %v", id, username, email, createdAt, updatedAt)
	// }
}

// CreateSettings handles POST requests to create new settings for a user
func CreateSettings(w http.ResponseWriter, r *http.Request) {
    var newSettings Settings
    _ = json.NewDecoder(r.Body).Decode(&newSettings)

    // Generate a temporary unique ID for the user
    newSettings.UserID = generateTempUserID()
    settingsMap[newSettings.UserID] = newSettings

    json.NewEncoder(w).Encode(newSettings)
}

// GetSettings handles GET requests to retrieve settings for a specific user
func GetSettings(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    userId := params["userId"]

    if settings, found := settingsMap[userId]; found {
        json.NewEncoder(w).Encode(settings)
    } else {
        w.WriteHeader(http.StatusNotFound)
    }
}

// generateTempUserID creates a simple temporary identifier for a new user
func generateTempUserID() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%v", rand.Int())
}
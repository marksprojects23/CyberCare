package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
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

var (
	db *sql.DB
)

// Temporary storage for settings
var settingsMap = make(map[string]Settings)


func main() {
	// PostgreSQL server credentials
	dbHost := "database-1.cl0i0y628wpg.us-east-1.rds.amazonaws.com"
	dbPort := 5432
	dbUser := "postgres"
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := "postgres"

// Connect to PostgreSQL database
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connString)
	if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal("Cannot connect to the database:", err)
    }

	r := mux.NewRouter()
    r.HandleFunc("/settings", CreateSettings).Methods("POST")
    r.HandleFunc("/settings/{userId}", GetSettings).Methods("GET")

    http.ListenAndServe(":8080", r)




	if err != nil {log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
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

    // Insert settings into the database
    var userId int
    err := db.QueryRow("INSERT INTO template (whitelist, blacklist) VALUES ($1, $2) RETURNING user_id", newSettings.Whitelist, newSettings.Blacklist).Scan(&userId)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    newSettings.UserID = strconv.Itoa(userId)
    json.NewEncoder(w).Encode(newSettings)
}


// GetSettings handles GET requests to retrieve settings for a specific user
func GetSettings(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    userId := params["userId"]

    var settings Settings
    settings.UserID = userId
    err := db.QueryRow("SELECT whitelist, blacklist FROM template WHERE user_id = $1", userId).Scan(&settings.Whitelist, &settings.Blacklist)
    if err == sql.ErrNoRows {
        http.NotFound(w, r)
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(settings)
}


// generateTempUserID creates a simple temporary identifier for a new user
func generateTempUserID() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%v", rand.Int())
}
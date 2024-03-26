package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
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
	ctx context.Context
	db sql.DB
)

// Temporary storage for settings
var settingsMap = make(map[string]Settings)

func main() {
	// SSH server credentials
	// sshHost := "ec2-52-86-177-235.compute-1.amazonaws.com"
	// sshPort := 22
	// sshUser := "ec2-user"
	// sshPrivateKeyPath := "bastion.pem"

	r := mux.NewRouter()
    r.HandleFunc("/settings", CreateSettings).Methods("POST")
    r.HandleFunc("/settings/{userId}", GetSettings).Methods("GET")

    http.ListenAndServe(":8080", r)

	// PostgreSQL server credentials
	dbHost := "localhost"
	dbPort := 5432
	dbUser := "postgres"
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := "postgres"

// Establish SSH tunnel
// sshClient, err := createSSHClient(sshUser, sshHost, sshPort, sshPrivateKeyPath)
// if err != nil {
// 	log.Fatalf("Failed to connect to SSH server: %v", err)
// }
// defer sshClient.Close()

// // Start local port forwarding
// localEndpoint := fmt.Sprintf("localhost:%d", dbPort)
// remoteEndpoint := fmt.Sprintf("%s:%d", dbHost, dbPort)
// err = forwardLocalPort(localEndpoint, remoteEndpoint, sshClient)
// if err != nil {
// 	log.Fatalf("Failed to forward local port: %v", err)
// }

// Connect to PostgreSQL database
connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", dbHost, dbPort, dbUser, dbPassword, dbName)
db, err := sql.Open("postgres", connString)
if err != nil {log.Fatalf("Failed to connect to PostgreSQL: %v", err)
}
rows, err := db.Query ("SELECT * FROM users;")
log.Output(1, fmt.Sprintf("%v", rows))
if err != nil {
	log.Fatalf("Failed to connect to PostgreSQL: %v", err)
}
defer db.Close()

for rows.Next() {
	var (
		id  int
		username string
		passwordHash string
		email string
		createdAt time.Time
		updatedAt time.Time

	)
	err := rows.Scan(&id, &username, &passwordHash, &email, &createdAt, &updatedAt)
	if err != nil {log.Fatalf("Error scanning row: ", err)
	}

	fmt.Printf("ID: %v, Password: %v, Email: %s, CreatedAt: %v, UpdatedAt: %v", id, username, email, createdAt, updatedAt)
}

// // Ping to check connection
// err = db.Ping()
// if err != nil {
// 	log.Fatalf("Cannot connect to the database: %v", err)
// }

// fmt.Println("Connected to PostgreSQL through SSH tunnel successfully!")
// }

// // createSSHClient creates an SSH client to connect to the SSH server
// func createSSHClient(user, host string, port int, privateKeyPath, caFilePath string) (*ssh.Client, error) {
// // Load private key
// keyBytes, err := os.ReadFile(privateKeyPath)
// if err != nil {
// 	return nil, fmt.Errorf("failed to read private key: %v", err)
// }
// signer, err := ssh.ParsePrivateKey(keyBytes)
// if err != nil {
// 	return nil, fmt.Errorf("failed to parse private key: %v", err)
// }

// // Create SSH client config
// config := &ssh.ClientConfig{
// 	User: user,
// 	Auth: []ssh.AuthMethod{
// 		ssh.PublicKeys(signer),
// 	},
// 	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
// }

// // Connect to SSH server
// conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
// if err != nil {
// 	return nil, fmt.Errorf("failed to connect to SSH server: %v", err)
// }

// return conn, nil
// }

// // forwardLocalPort forwards local port to remote endpoint via SSH tunnel
// func forwardLocalPort(localEndpoint, remoteEndpoint string, sshClient *ssh.Client) error {
// // Listen on local port
// listener, err := net.Listen("tcp", localEndpoint)
// if err != nil {
// 	return fmt.Errorf("failed to start local listener: %v", err)
// }
// defer listener.Close()

// // Handle incoming connections
// for {
// 	localConn, err := listener.Accept()
// 	if err != nil {
// 		return fmt.Errorf("failed to accept connection: %v", err)
// 	}

// 	// Connect to remote endpoint via SSH tunnel
// 	remoteConn, err := sshClient.Dial("tcp", remoteEndpoint)
// 	if err != nil {
// 		localConn.Close()
// 		return fmt.Errorf("failed to establish SSH tunnel: %v", err)
// 	}

// 	// Forward traffic between local and remote connections
// 	go func() {
// 		defer localConn.Close()
// 		defer remoteConn.Close()

// 		// Copy data from local to remote
// 		go func() {
// 			if _, err := io.Copy(remoteConn, localConn); err != nil {
// 				log.Printf("error copying local to remote: %v", err)
// 			}
// 		}()

// 		// Copy data from remote to local
// 		if _, err := io.Copy(localConn, remoteConn); err != nil {
// 			log.Printf("error copying remote to local: %v", err)
// 		}
// 	}()
// }
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
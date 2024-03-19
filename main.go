// // main.go
// package main

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// const (
// 	mongoURI      = "mongodb://localhost:27017"
// 	databaseName  = "cybercare"
// 	collectionName = "users"
// )

// var client *mongo.Client

// type User struct {
// 	Username string `json:"username" bson:"username"`
// 	Password string `json:"password" bson:"password"`
// }

// func main() {
// 	// Initialize MongoDB client
// 	clientOptions := options.Client().ApplyURI(mongoURI)
// 	var err error
// 	client, err = mongo.NewClient(clientOptions)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Disconnect(ctx)

// 	// Initialize Gin router
// 	router := gin.Default()

// 	// Define routes
// 	router.POST("/register", registerHandler)
// 	router.POST("/login", loginHandler)

// 	// Run the server
// 	err = router.Run(":8080")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func registerHandler(c *gin.Context) {
// 	var newUser User
// 	if err := c.ShouldBindJSON(&newUser); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Check if the username already exists
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	collection := client.Database(databaseName).Collection(collectionName)
// 	existingUser := User{}
// 	err := collection.FindOne(ctx, bson.M{"username": newUser.Username}).Decode(&existingUser)
// 	if err == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
// 		return
// 	}

// 	// Hash the password (You should use a proper hashing library in a production environment)
// 	// For simplicity, we are storing the plain text password here.
// 	// In production, you should always hash passwords before storing them.
// 	newUser.Password = newUser.Password

// 	// Insert the new user into the database
// 	_, err = collection.InsertOne(ctx, newUser)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
// }

// func loginHandler(c *gin.Context) {
// 	var loginRequest User
// 	if err := c.ShouldBindJSON(&loginRequest); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Check if the user exists
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	collection := client.Database(databaseName).Collection(collectionName)
// 	existingUser := User{}
// 	err := collection.FindOne(ctx, bson.M{"username": loginRequest.Username}).Decode(&existingUser)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 		return
// 	}

// 	// Verify the password (You should use a proper authentication library in a production environment)
// 	// For simplicity, we are comparing plain text passwords here.
// 	// In production, you should always verify passwords using hashed values.
// 	if existingUser.Password != loginRequest.Password {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
// }

package main

import (
	"fmt"
	"net/http"
)

// this is for Mac users only
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, you've reached the home page!")
    })

    fmt.Println("Server starting on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        panic(err)
    }
}


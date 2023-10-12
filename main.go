package main

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var SECRET_KEY = []byte("gosecretkey")

type User struct {
	Name     string `json:"name" bson:"name"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

var client *mongo.Client

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}

func userSignup(response http.ResponseWriter, request *http.Request) {
	//response.Header().Set("Content-Type", "application/json")
	//response.Header().Set("Access-Control-Allow-Origin", "*")
	//response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	//response.Header().Set("Access-Control-Allow-Headers", "*")
	enableCors(&response)
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	user.Password = getHash([]byte(user.Password))
	collection := client.Database("GODB").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

func userLogin(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	var user User
	var dbUser User
	json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("GODB").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbUser)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	userPass := []byte(user.Password)
	dbPass := []byte(dbUser.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		log.Println(passErr)
		response.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}
	jwtToken, err := GenerateJWT()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}
	response.Write([]byte(`{"token":"` + jwtToken + `"}`))

}
func allusers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	user.Password = getHash([]byte(user.Password))
	collection := client.Database("GODB").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var documents []User

	// Iterate through the cursor and decode each document
	for cursor.Next(context.TODO()) {
		var doc User
		err := cursor.Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}
		documents = append(documents, doc)
	}

	// Check for any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	// Convert the documents slice to JSON
	jsonBytes, err := json.Marshal(documents)
	if err != nil {
		log.Fatal(err)
	}

	// Set the response content type and write the JSON response
	response.Header().Set("Content-Type", "application/json")
	response.Write(jsonBytes)

}
func main() {
	log.Println("Starting the application")

	router := mux.NewRouter()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	router.HandleFunc("/api/user/login", userLogin).Methods("POST")
	router.HandleFunc("/api/user/signup", userSignup).Methods("POST")
	router.HandleFunc("/users", allusers).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))

}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

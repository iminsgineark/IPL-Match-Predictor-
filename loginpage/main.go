package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	client   *mongo.Client
	userColl *mongo.Collection
)

func main() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		log.Fatal("Error creating MongoDB client:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	userColl = client.Database("authDB").Collection("users")
	fmt.Println("Connected to MongoDB!")

	http.HandleFunc("/", LoginPage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/signup", SignupPage)
	http.HandleFunc("/welcome", WelcomePage)

	fmt.Println("Server started on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if password != confirmPassword {
			fmt.Fprintf(w, "Passwords do not match. Please try again.")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var existingUser bson.M
		err := userColl.FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)
		if err == nil {
			fmt.Fprintf(w, "Username already exists. Please choose another username.")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password.", http.StatusInternalServerError)
			log.Println("Error hashing password:", err)
			return
		}

		_, err = userColl.InsertOne(ctx, bson.M{
			"username": username,
			"email":    email,
			"password": string(hashedPassword),
		})
		if err != nil {
			http.Error(w, "Failed to save user data.", http.StatusInternalServerError)
			log.Println("Error inserting user:", err)
			return
		}

		fmt.Printf("New user signup: Username - %s, Email - %s\n", username, email)

		http.Redirect(w, r, "http://localhost/model/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user bson.M
		err := userColl.FindOne(ctx, bson.M{"username": username}).Decode(&user)
		if err != nil {
			fmt.Fprintf(w, "Invalid username or password. Please try again.")
			return
		}

		storedPassword := user["password"].(string)

		if strings.HasPrefix(storedPassword, "$2a$") || strings.HasPrefix(storedPassword, "$2b$") || strings.HasPrefix(storedPassword, "$2y$") {
			err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
			if err != nil {
				fmt.Fprintf(w, "Invalid username or password. Please try again.")
				return
			}
		} else {
			if storedPassword != password {
				fmt.Fprintf(w, "Invalid username or password. Please try again.")
				return
			}
		}

		http.Redirect(w, r, "http://localhost/model/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func WelcomePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost/model/", http.StatusSeeOther)
}

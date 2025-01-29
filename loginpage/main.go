package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal("Error disconnecting MongoDB:", err)
		}
	}()

	userColl = client.Database("authDB").Collection("users")
	fmt.Println("Connected to MongoDB!")

	// Serve static files (CSS, images, etc.)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handlers for login and signup pages
	http.HandleFunc("/", LoginPage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/signup", SignupPage)
	http.HandleFunc("/model/", model)

	server := &http.Server{Addr: ":8080"}
	go func() {
		fmt.Println("Server started on http://0.0.0.0:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	fmt.Println("\nShutting down server...")
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	fmt.Println("Server exited gracefully.")
}

func renderTemplate(w http.ResponseWriter, templateFile string, data interface{}) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		log.Println("Template Error:", err)
		return
	}
	tmpl.Execute(w, data)
}

func handleErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
	log.Println("Error:", message)
}

func SignupPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if password != confirmPassword {
			handleErrorResponse(w, "Passwords do not match.", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var existingUser bson.M
		err := userColl.FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)
		if err == nil {
			handleErrorResponse(w, "Username already exists. Please choose another.", http.StatusConflict)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			handleErrorResponse(w, "Failed to hash password.", http.StatusInternalServerError)
			return
		}

		_, err = userColl.InsertOne(ctx, bson.M{
			"username": username,
			"email":    email,
			"password": string(hashedPassword),
		})
		if err != nil {
			handleErrorResponse(w, "Failed to save user data.", http.StatusInternalServerError)
			return
		}

		log.Printf("New user signup: Username - %s, Email - %s\n", username, email)
		http.Redirect(w, r, "/model", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "templates/signup.html", nil)
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
			handleErrorResponse(w, "Invalid username or password.", http.StatusUnauthorized)
			return
		}

		storedPassword := user["password"].(string)
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		if err != nil {
			handleErrorResponse(w, "Invalid username or password.", http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/model", http.StatusSeeOther)
		return
	}

	renderTemplate(w, "templates/login.html", nil)
}

func model(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the model page!"))
}

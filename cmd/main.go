package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/alexski/grocer-backend/handler"

	"github.com/joho/godotenv"
)

// Current version of application
var Version = "0.0.1"

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	a := App{}
	a.Initialize(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	// a.Run(":8010")

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "API Connected")
	})
	router.HandleFunc("/user", handler.GetAllUsers)
	router.HandleFunc("/user/", handler.GetSingleUser)

	port := ":8080"
	_, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("Defaulting port.")
	} else {
		port = ":" + os.Getenv("PORT")
	}

	fmt.Println("API Version: " + Version)
	fmt.Println("Running server on port " + port)
	log.Fatal(http.ListenAndServe(port, router))
}

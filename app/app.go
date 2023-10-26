package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// This struct exposes references to the router and the database that the application uses.
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// To be useful and testable, App will need two methods that initialize and run the application.
// Takes in the details required to connect to the database.
// It will create a database connection and wire up the routes
func (a *App) Initialize(user string, password string, dbname string, host string, port string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
}

// Simply starts the application
func (a *App) Run(addr string) {}

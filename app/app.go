package app

import (
	"database/sql"

	"github.com/gorilla/mux"
)

// This struct exposes references to the router and the database that the application uses.
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// To be useful and testable, App will need two methods that initialize and run the application.
// Takes in the details required to connect to the database.
// It will create a database connection and wire up the routes
func (a *App) Initialize(user, password, dbname string) {}

// Simply starts the application
func (a *App) Run(addr string) {}

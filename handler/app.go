package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"grocer-backend/model"

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
	a.initializeRoutes()
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	u := model.User{ID: id}
	if err := u.getUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) initializeRoutes() {
	// a.Router.HandleFunc("/users", a.getUsers).Methods("GET")
	// a.Router.HandleFunc("/user", createUser).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", a.getUser).Methods("GET")
	// a.Router.HandleFunc("/user/{id:[0-9]+}", a.updateUser).Methods("PUT")
	// a.Router.HandleFunc("/user/{id:[0-9]+}", a.deleteUser).Methods("DELETE")
}

// Simply starts the application
func (a *App) Run(addr string) {}

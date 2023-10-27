package grocer_tests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"grocer-backend/handler"

	"github.com/joho/godotenv"
)

var a handler.App

func TestMain(m *testing.M) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	a.Initialize(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"))

	ensureTableExists()
	code := m.Run()
	clearTables()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(userTableCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(recipeTableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTables() {
	clearUserTable()
	clearRecipeTable()
}

func TestEmptyTable(t *testing.T) {
	clearTables()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

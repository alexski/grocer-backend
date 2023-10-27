package grocer_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

const userTableCreationQuery = `
CREATE SEQUENCE IF NOT EXISTS public.users_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

CREATE TABLE IF NOT EXISTS public.users
(
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    username character varying(100) COLLATE pg_catalog."default" NOT NULL,
    password_hash character varying(100) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

ALTER SEQUENCE public.users_id_seq
	OWNED BY users.id;
`

func clearUserTable() {
	a.DB.Exec("DELETE FROM users")
	a.DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

func TestGetNonExistentUser(t *testing.T) {
	clearUserTable()

	req, _ := http.NewRequest("GET", "/user/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "User not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
	}
}

func TestCreateUser(t *testing.T) {

	clearUserTable()

	var jsonStr = []byte(`{"username":"test_user", "password": "test_hash"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["username"] != "test_user" {
		t.Errorf("Expected username name to be 'test_user'. Got '%v'", m["username"])
	}

	if m["password"] != "test_hash" {
		t.Errorf("Expected user's password to be 'test_hash'. Got '%v'", m["password"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected user ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetUser(t *testing.T) {
	clearUserTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateUser(t *testing.T) {

	clearUserTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)
	var ogUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &ogUser)

	var jsonStr = []byte(`{"username":"test_user_updated", "password": "test_hash_updated"}`)
	req, _ = http.NewRequest("PUT", "/user/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != ogUser["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", ogUser["id"], m["id"])
	}

	if m["username"] == ogUser["username"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", ogUser["username"], m["username"], m["username"])
	}

	if m["password"] == ogUser["password"] {
		t.Errorf("Expected the password hash to change from '%v' to '%v'. Got '%v'", ogUser["password"], m["password"], m["password"])
	}
}

func TestDeleteUser(t *testing.T) {
	clearUserTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/user/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/user/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func addUsers(count int) {
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO users(username, password_hash) VALUES($1, $2)", "test_user_"+strconv.Itoa(i), "test_hash")
	}
}

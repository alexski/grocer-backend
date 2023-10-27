package grocer_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
)

const recipeTableCreationQuery = `
CREATE SEQUENCE IF NOT EXISTS public.recipes_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 2147483647
    CACHE 1;

CREATE TABLE IF NOT EXISTS public.recipes
(
    id integer NOT NULL DEFAULT nextval('recipes_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    image_filename character varying(500) COLLATE pg_catalog."default",
    filename character varying(500) COLLATE pg_catalog."default" NOT NULL,
    type character varying(100) COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT recipes_pkey PRIMARY KEY (id)
);

ALTER SEQUENCE public.recipes_id_seq
	OWNED BY recipes.id;
`

func clearRecipeTable() {
	a.DB.Exec("DELETE FROM recipes")
	a.DB.Exec("ALTER SEQUENCE recipes_id_seq RESTART WITH 1")
}

func TestGetNonExistentRecipe(t *testing.T) {
	clearRecipeTable()

	req, _ := http.NewRequest("GET", "/recipe/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Recipe not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Recipe not found'. Got '%s'", m["error"])
	}
}

func TestCreateRecipe(t *testing.T) {

	clearRecipeTable()

	var jsonStr = []byte(`{"title":"chicken recipe", "image": "chicken.png", "filename": "chicken.md", "type": "entree"}`)
	req, _ := http.NewRequest("POST", "/recipe", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "chicken recipe" {
		t.Errorf("Expected recipe title to be 'chicken recipe'. Got '%v'", m["title"])
	}

	if m["image"] != "chicken.png" {
		t.Errorf("Expected recipe's password to be 'chicken.png'. Got '%v'", m["image"])
	}

	if m["filename"] != "chicken.md" {
		t.Errorf("Expected recipe's filename to be 'chicken.md'. Got '%v'", m["filename"])
	}

	if m["type"] != "entree" {
		t.Errorf("Expected recipe's type to be 'entree'. Got '%v'", m["type"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected recipe ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetRecipe(t *testing.T) {
	clearRecipeTable()
	addRecipes(1)

	req, _ := http.NewRequest("GET", "/recipe/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateRecipe(t *testing.T) {
	clearRecipeTable()
	addRecipes(1)

	req, _ := http.NewRequest("GET", "/recipe/1", nil)
	response := executeRequest(req)
	var ogRecipe map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &ogRecipe)

	var jsonStr = []byte(`{"title":"chicken updated", "image": "chicken_modded.png", "filename": "chicken_modded.md", "type": "side"}`)
	req, _ = http.NewRequest("PUT", "/recipe/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != ogRecipe["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", ogRecipe["id"], m["id"])
	}

	if m["title"] == ogRecipe["title"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", ogRecipe["name"], m["name"], m["name"])
	}

	if m["image"] == ogRecipe["image"] {
		t.Errorf("Expected the recipe's image filename to change from '%v' to '%v'. Got '%v'", ogRecipe["image"], m["image"], m["image"])
	}

	if m["filename"] == ogRecipe["filename"] {
		t.Errorf("Expected the recipe's filename to change from '%v' to '%v'. Got '%v'", ogRecipe["filename"], m["filename"], m["filename"])
	}

	if m["type"] == ogRecipe["type"] {
		t.Errorf("Expected the recipe's type to change from '%v' to '%v'. Got '%v'", ogRecipe["type"], m["type"], m["type"])
	}
}

func TestDeleteRecipe(t *testing.T) {
	clearRecipeTable()
	addRecipes(1)

	req, _ := http.NewRequest("GET", "/recipe/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/recipe/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/recipe/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func addRecipes(count int) {
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO recipes(name, image_filename, filename, type) VALUES($1, $2, $3, $4)", "chicken recipe "+strconv.Itoa(i), "chicken.png", "chicken.md", "entree")
		a.DB.Exec("INSERT INTO recipes(name, image_filename, filename, type) VALUES($1, $2, $3, $4)", "steak recipe "+strconv.Itoa(i), "steak.png", "steak.md", "entree")
	}
}

package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func GetSingleUser(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	if len(p) < 3 {
		http.Error(w, "No Id Provided", http.StatusBadRequest)
		return
	} else if len(p) > 1 {
		code, err := strconv.Atoi(p[2])
		if err == nil {
			fmt.Fprint(w, strconv.Itoa(code))
			return
		} else {
			http.Error(w, "Id provided not a number", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "No Id Provided", http.StatusBadRequest)
		return
	}
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprintf(w, "These are all the users.")
}

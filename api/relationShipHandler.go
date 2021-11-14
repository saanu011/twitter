package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"twitter/db"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var followUser struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&followUser); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// check if both userIDs are different
	if followUser.UserID == userID {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("you can't follow yourself"))
		return
	}
	txo, failed := db.BeginTx(database)
	if failed {
		return
	}

	// create relation
	relation := db.Relationship{
		FollowedID: followUser.UserID,
		FollowerID: userID,
	}
	err = relation.Create(txo)
	if db.TerminateTx(txo, err) {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating relationship: %v", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

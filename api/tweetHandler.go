package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"twitter/db"
)

func PostTweet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var createReq db.Tweet
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	txo, failed := db.BeginTx(database)
	if failed {
		return
	}

	// create tweet
	tweet := db.Tweet{
		Content: createReq.Content,
		UserID:  userID,
	}
	err = tweet.Create(txo)
	if db.TerminateTx(txo, err) {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating authClient: %v", err))
		return
	}

	if err := json.NewEncoder(w).Encode(tweet); err != nil {
		logrus.WithError(err).Warn("Error encoding JSON response")
	}
}

func GetUserTweets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	// Begin transaction
	txo, failed := db.BeginTx(database)
	if failed {
		return
	}

	user, err := db.GetUserTweets(txo, userID)
	if db.TerminateTx(txo, err) {
		if strings.Contains(err.Error(), "no rows found") {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("no user found with ID: %d", userID))
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error retrieving user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func GetAllFollowedUsersTweets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	// Begin transaction
	txo, failed := db.BeginTx(database)
	if failed {
		return
	}

	tweetFeed, err := db.GetFollowedUserTweets(txo, userID)
	if db.TerminateTx(txo, err) {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(tweetFeed) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	respondWithJSON(w, http.StatusOK, tweetFeed)
}

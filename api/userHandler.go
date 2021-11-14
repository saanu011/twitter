package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
	"net/http"
	"strconv"
	"strings"
	"twitter/db"
	"twitter/utils"
)

type loginInfo struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// userCreateReq data request
type userCreateReq struct {
	Name            string      `json:"name" validate:"required"`
	Email           string      `json:"email" validate:"required"`
	Address         null.String `json:"address" db:"address"`
	Password        string      `json:"password" validate:"required,min=8"`
	PasswordConfirm string      `json:"password_confirm" validate:"required,min=8"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	var credentials loginInfo
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	txo, failed := db.BeginTx(database)
	if failed {
		return
	}
	// check if credentials exists or not and get its auth info by email
	authClient, err := db.GetAuthClientByClientRef(txo, credentials.Email)
	if db.TerminateTxIfError(txo, err) {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("no user found with provided credentials: %v", err))
		return
	}

	// match password hash
	if !utils.HashMatchesPassword(authClient.ClientSecret, credentials.Password) {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("no user found with provided credentials: %v", err))
		return
	}

	// if everything's okay, go ahead and create access token
	accessToken, err := db.CreateAccessToken(txo, authClient)
	if db.TerminateTx(txo, err) {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}
	response := struct {
		Token string `json:"token"`
	}{
		Token: accessToken.AccessToken,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.WithError(err).Warn("Error encoding JSON response")
	}
}

func UserCreate(w http.ResponseWriter, r *http.Request) {

	var createReq userCreateReq
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	txo, failed := db.BeginTx(database)
	if failed {
		return
	}

	if createReq.Password != createReq.PasswordConfirm {
		respondWithError(w, http.StatusBadRequest, "password does not match")
		return
	}
	// create auth client
	authClient := db.AuthClient{
		ClientRef:    createReq.Email,
		ClientSecret: utils.Hash(createReq.Password),
	}
	// create user
	user := db.User{
		Name:    createReq.Name,
		Email:   createReq.Email,
		Address: createReq.Address,
	}
	err := user.Create(txo)
	if db.TerminateTxIfError(txo, err) {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating user: %v", err))
		return
	}

	err = authClient.Create(txo)
	if db.TerminateTx(txo, err) {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error creating authClient: %v", err))
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		logrus.WithError(err).Warn("Error encoding JSON response")
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := db.GetUserById(txo, userID)
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

func GetUsers(w http.ResponseWriter, r *http.Request) {

	// Begin transaction
	txo, failed := db.BeginTx(database)
	if failed {
		return
	}

	user, err := db.GetUsers(txo)
	if db.TerminateTx(txo, err) {
		respondWithError(w, http.StatusInternalServerError, "error retrieving user")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

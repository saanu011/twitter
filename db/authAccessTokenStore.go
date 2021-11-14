package db

import (
	"bytes"
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"os"
	"strconv"
	"time"
)

func (token *AuthAccessToken) Upsert(txo *TxO) error {

	var q bytes.Buffer
	q.WriteString(fmt.Sprintf(`
		INSERT INTO auth_access_token
		(
			auth_client_id,
			access_token,
			expires_in,
			expires_on
		)
		VALUES `,
	))

	q.WriteString(`(?,?,?,?)`)
	parameters := []interface{}{
		token.AuthClientID,
		token.AccessToken,
		token.ExpiresIn,
		token.ExpiresOn,
	}

	q.WriteString(` ON DUPLICATE KEY UPDATE access_token = VALUES(access_token), expires_in = VALUES(expires_in), expires_on = VALUES(expires_on) `)

	r, err := txo.Exec(q.String(), parameters...)
	if err != nil {
		return fmt.Errorf("")
	}

	lastID, err := r.LastInsertId()
	if err != nil {
		return fmt.Errorf("")
	}

	token.ID = int(lastID)

	return nil
}

// CreateAccessToken ...
func CreateAccessToken(txo *TxO, authClient *AuthClient) (accessToken AuthAccessToken, err error) {

	ttl, err := strconv.Atoi(os.Getenv("TTL"))
	if err != nil {
		return accessToken, err
	}

	user, err := GetUserByEmail(txo, authClient.ClientRef)
	if err != nil {
		return accessToken, err
	}

	accessTokenStr, err := GenerateToken(user)
	if err != nil {
		return accessToken, err
	}

	accessToken = AuthAccessToken{
		AuthClientID: authClient.ID,
		AccessToken:  accessTokenStr,
		ExpiresIn:    os.Getenv("TTL"),
		ExpiresOn:    time.Now().Add(time.Duration(ttl)),
	}

	// update a record for access token
	err = accessToken.Upsert(txo)
	if err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

// GenerateToken ...
func GenerateToken(u *User) (string, error) {

	ttl, err := strconv.Atoi(os.Getenv("TTL"))
	if err != nil {
		return "", err
	}
	return jwt.NewWithClaims(jwt.GetSigningMethod(os.Getenv("SIGNING_ALGORITHM")), jwt.MapClaims{
		"id":  u.UserID,
		"n":   u.Name,
		"e":   u.Email,
		"exp": time.Now().Add(time.Duration(ttl) * time.Minute).Unix(),
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
}

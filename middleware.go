package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/form3tech-oss/jwt-go"
)

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

		if len(secret) == 0 {
			return nil, errors.New("JWT Secret Not Set")
		}

		return []byte(secret), nil
	},

	Extractor: extractTokenFromHeader,
	Debug:     false,
})

func extractTokenFromHeader(r *http.Request) (string, error) {

	// get Authorization Header
	auth := r.Header.Get("Authorization")

	// Authorization token:
	// Bearer XXXXX
	//
	// Get actual token
	splitAuth := strings.Split(auth, " ")
	if len(splitAuth) < 2 {
		return "", errors.New("Invalid Auth: length of split authorization token is less than 2")
	}

	return splitAuth[1], nil
}

// addCustomContext to the request so that handlers can use helper fields and methods
func addCustomContext(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// new context with trace ID as UUID
		ctx := FromContext(r.Context())

		ctx = SetHeaderWithUserInfo(ctx, r)

		// copy to the given request pointer
		*r = *r.WithContext(ctx)
		// pass it along for the next one in the chain
		inner.ServeHTTP(w, r)
	})
}

type Context struct {
	// Context is the base underlying context.
	context.Context
	// Log should be set with critical log fields that will be passed along
	// and logged downstream. Set this explicitly to define a custom log entry,
	// but setting will lose any previously set fields.
	Log *logrus.Entry
}

// FromContext creates a new Context, using parent as the base context.Context.  If
// parent is already contexture.Context then a new Context will be returned as contexture.Context,
// preserving any existing configuration, but in a new instance.
func FromContext(parent context.Context) Context {
	// preserve context when given
	if ctx, isContext := parent.(Context); isContext {
		return ctx
	}

	return Context{
		Context: parent,
		Log:     logrus.NewEntry(logrus.StandardLogger()),
	}
}

// SetHeaderWithUserInfo ...
func SetHeaderWithUserInfo(ctx Context, r *http.Request) Context {
	claims, _ := mapToClaims(r)

	id, _ := claims["_id"].(string)
	name, _ := claims["name"].(string)
	email, _ := claims["email"].(string)

	ctx.Context = context.WithValue(ctx, "_id", id)
	ctx.Context = context.WithValue(ctx, "name", name)
	ctx.Context = context.WithValue(ctx, "email", email)
	return ctx
}

func mapToClaims(r *http.Request) (jwt.MapClaims, error) {
	// get user context from request. This should be set by another process, usually the JWT middleware, or this will fail.
	user := r.Context().Value("user")

	// convert user token to JWToken
	jwtToken, ok := user.(*jwt.Token)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("user context not of type *jwt.Token")
	}

	// get the map of claims from user context
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("could not get claims from JWT")
	}
	return claims, nil
}

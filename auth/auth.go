package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sanderhahn/htpasswd/db"
	jwt "github.com/sanderhahn/htpasswd/jwt"
	"github.com/sanderhahn/htpasswd/wrap"
)

func parseSecret(secret string) (s jwt.Secret, err error) {
	err = json.Unmarshal([]byte(secret), &s)
	return s, err
}

// Input struct
type Input struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Output struct
type Output struct {
	Token   string `json:"token"`
	Expires int64  `json:"expires"`
}

// Authenticate a user against the database
func Authenticate(w wrap.ResponseWriter, r *wrap.Request) {
	var secret jwt.Secret
	if err := r.ParseHeader("hasura-graphql-jwt-secret", &secret); err != nil {
		w.Error(fmt.Errorf("syntax error in hasura-graphql-jwt-secret"), http.StatusInternalServerError)
		return
	}

	var d db.Database
	if r.HasHeader("x-htpasswd") {
		if err := r.ParseHeader("x-htpasswd", &d); err != nil {
			w.Error(fmt.Errorf("syntax error in x-htpasswd"), http.StatusInternalServerError)
			return
		}
	} else {
		d = db.DefaultDatabase()
	}

	var input Input
	if err := r.ParseInput(&input); err != nil {
		w.Error(err, http.StatusBadRequest)
		return
	}

	user := d.Authenticate(input.Username, input.Password)
	if user == nil {
		w.Error(fmt.Errorf("Forbidden"), http.StatusForbidden)
		return
	}

	expiresIn := time.Hour * time.Duration(24)
	token, err := secret.Sign(expiresIn, user.Username, user.Role)
	if err != nil {
		w.Error(err, http.StatusInternalServerError)
		return
	}

	w.JSON(Output{
		Token:   token,
		Expires: int64(expiresIn.Seconds()),
	})
}

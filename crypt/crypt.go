package crypt

import (
	"fmt"
	"net/http"

	"github.com/sanderhahn/htpasswd/wrap"
	"golang.org/x/crypto/bcrypt"
)

// Input is input
type Input struct {
	Password string `json:"password"`
}

// Output is output
type Output struct {
	Hashed string `json:"hashed"`
}

// Crypt a password
func Crypt(w wrap.ResponseWriter, r *wrap.Request) {
	var input Input

	if err := r.ParseInput(&input); err != nil {
		w.Error(err, http.StatusBadRequest)
		return
	}

	if input.Password == "" {
		w.Error(fmt.Errorf("password is empty"), http.StatusBadRequest)
		return
	}

	in := []byte(input.Password)
	pw, err := bcrypt.GenerateFromPassword(in, bcrypt.MinCost)
	if err != nil {
		w.Error(err, http.StatusInternalServerError)
		return
	}

	w.JSON(Output{
		Hashed: string(pw),
	})
	return
}

// CheckCrypt compares the hashed and password and returns nil on success
func CheckCrypt(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

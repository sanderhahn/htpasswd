package crypt

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sanderhahn/htpasswd/wrap"
	"golang.org/x/crypto/bcrypt"
)

// Input struct
type Input struct {
	Password string `json:"password"`
}

// Output struct
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

	w.JSON(Output{
		Hashed: Password(input.Password),
	})
	return
}

// Password hashes a password
func Password(password string) string {
	buf, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

// Check compares the hashed and password and returns nil on success
func Check(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

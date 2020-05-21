package db

import (
	"encoding/json"

	"github.com/sanderhahn/htpasswd/crypt"
)

// User represents a user
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// Database is an array of user entries
type Database []User

// DefaultDatabase has a static database of user entries
func DefaultDatabase() Database {
	secret := crypt.Password("secret")
	return Database{
		{
			Username: "user",
			Password: secret,
			Role:     "user",
		},
		{
			Username: "manager",
			Password: secret,
			Role:     "manager",
		},
		{
			Username: "admin",
			Password: secret,
			Role:     "admin",
		},
	}
}

// ParseDatabase loads database from json
func ParseDatabase(data string) (Database, error) {
	var d Database
	err := json.Unmarshal([]byte(data), &d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// LookupUser by username
func (d *Database) LookupUser(username string) *User {
	for _, user := range *d {
		if user.Username == username {
			return &user
		}
	}
	return nil
}

// Authenticate checks for the correct authentication
func (d *Database) Authenticate(username, password string) *User {
	user := d.LookupUser(username)
	if user != nil && crypt.Check(user.Password, password) == nil {
		return user
	}
	return nil
}

package whoami

import (
	"fmt"
	"net/http"

	"github.com/sanderhahn/htpasswd/wrap"
)

// Input struct
type Input struct {
}

// Output struct
type Output struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Whoami shows the username/role based on the session variables
func Whoami(w wrap.ResponseWriter, r *wrap.Request) {
	username := r.Payload.SessionVariables["x-hasura-user-id"]
	role := r.Payload.SessionVariables["x-hasura-role"]
	if username == "" && role != "admin" {
		w.Error(fmt.Errorf("Forbidden"), http.StatusForbidden)
		return
	}

	w.JSON(Output{
		Username: username,
		Role:     role,
	})
}

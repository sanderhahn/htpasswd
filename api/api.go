package api

import (
	"fmt"
	"net/http"

	"github.com/sanderhahn/htpasswd/auth"
	"github.com/sanderhahn/htpasswd/crypt"
	"github.com/sanderhahn/htpasswd/whoami"
	"github.com/sanderhahn/htpasswd/wrap"
)

// DispatchHandler dispatches api requests
func DispatchHandler(_w http.ResponseWriter, _r *http.Request) {
	w := wrap.NewResponseWriter(_w)
	r, err := wrap.NewRequest(_r)
	if err != nil {
		w.Error(err, http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		w.Error(fmt.Errorf("Invalid method"), http.StatusMethodNotAllowed)
		return
	}

	switch r.Payload.Action.Name {
	case "crypt":
		crypt.Crypt(w, r)

	case "authenticate":
		auth.Authenticate(w, r)

	case "whoami":
		whoami.Whoami(w, r)

	default:
		w.Error(fmt.Errorf("Invalid action %s", r.Payload.Action.Name), http.StatusBadRequest)
	}
}

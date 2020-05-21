package whoami

import (
	"testing"

	"github.com/sanderhahn/htpasswd/wrap"
	"github.com/sanderhahn/htpasswd/wraptest"
)

func TestWhoamiForbidden(t *testing.T) {
	w := &wraptest.MockResponseWriter{}
	r := wraptest.NewRequest(t, "whoami", nil, nil)

	Whoami(w, r)

	if w.ErrorOutput.Error() != "Forbidden" {
		t.Fail()
	}
}

func TestWhoami(t *testing.T) {
	w := &wraptest.MockResponseWriter{}
	r := wraptest.NewRequest(t, "whoami", Input{}, wrap.SessionVariables{
		"x-hasura-user-id": "admin",
		"x-hasura-role":    "admin",
	})

	Whoami(w, r)

	out := w.Output.(Output)
	if out.Username != "admin" || out.Role != "admin" {
		t.Fail()
	}
}

func TestWhoamiAdmin(t *testing.T) {
	w := &wraptest.MockResponseWriter{}
	r := wraptest.NewRequest(t, "whoami", Input{}, wrap.SessionVariables{
		"x-hasura-role": "admin",
	})

	Whoami(w, r)

	out := w.Output.(Output)
	if out.Role != "admin" {
		t.Fail()
	}
}

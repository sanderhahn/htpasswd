package crypt

import (
	"strings"
	"testing"

	"github.com/sanderhahn/htpasswd/wraptest"
)

func TestCrypt(t *testing.T) {
	w := &wraptest.MockResponseWriter{}
	r := wraptest.NewRequest(t, "crypt", Input{
		Password: "test",
	}, nil)

	Crypt(w, r)

	out := w.Output.(Output)
	if !strings.HasPrefix(out.Hashed, "$2a$04$") {
		t.Fail()
	}
}

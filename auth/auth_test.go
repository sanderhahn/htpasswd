package auth

import (
	"testing"

	jwt "github.com/sanderhahn/htpasswd/jwt"
	"github.com/sanderhahn/htpasswd/wraptest"
)

func TestNoSecretHeader(t *testing.T) {
	w := &wraptest.MockResponseWriter{}
	r := wraptest.NewRequest(t, "authenticate", Input{
		Username: "admin",
		Password: "secret",
	}, nil)

	Authenticate(w, r)

	if w.ErrorOutput.Error() != "Invalid JWT secret" {
		t.Fail()
	}
}

func fakeSecret(t *testing.T) string {
	return string(wraptest.MustMarshal(t, jwt.Secret{
		Type: "HS256",
		Key:  "secret",
	}))
}

func TestWhoami(t *testing.T) {
	w := &wraptest.MockResponseWriter{}
	r := wraptest.NewRequest(t, "authenticate", Input{
		Username: "admin",
		Password: "secret",
	}, nil)
	r.Header.Add("hasura-graphql-jwt-secret", fakeSecret(t))

	Authenticate(w, r)

	out := w.Output.(Output)
	if len(out.Token) == 0 {
		t.Fail()
	}
}

func TestWhoamiWrong(t *testing.T) {
	w := &wraptest.MockResponseWriter{}
	r := wraptest.NewRequest(t, "authenticate", Input{
		Username: "admin",
		Password: "wrong",
	}, nil)
	r.Header.Add("hasura-graphql-jwt-secret", fakeSecret(t))

	Authenticate(w, r)

	if w.ErrorOutput.Error() != "Forbidden" {
		t.Fail()
	}
}

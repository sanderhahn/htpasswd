package db

import "testing"

func TestDB(t *testing.T) {
	d := DefaultDatabase()
	if d.Authenticate("admin", "secret") == nil {
		t.Fail()
	}
}

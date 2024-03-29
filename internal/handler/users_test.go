package handler

import "testing"


func TestUsersHandler(t *testing.T) {

	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	}

	t.Run("saying hello to people", func(t *testing.T) {
		// got := Hello("Chris")
		// want := "Hello, Chris"
		// assertCorrectMessage(t, got, want)
	})

	//t.Run("empty string defaults to 'world'", func(t *testing.T) {
	//	got := Hello("")
	//	want := "Hello, World"
	//	assertCorrectMessage(t, got, want)
	//})
}

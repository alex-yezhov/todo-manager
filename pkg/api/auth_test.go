package api

import "testing"

func TestMakeTokenAndValidateToken(t *testing.T) {
	password := "12345"

	token, err := makeToken(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token == "" {
		t.Fatal("token is empty")
	}

	if !validateToken(token, password) {
		t.Fatal("token should be valid")
	}
}

func TestValidateToken_WrongPassword(t *testing.T) {
	password := "12345"

	token, err := makeToken(password)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if validateToken(token, "54321") {
		t.Fatal("token should be invalid for different password")
	}
}

func TestValidateToken_BadFormat(t *testing.T) {
	if validateToken("not-a-token", "12345") {
		t.Fatal("bad token format must be invalid")
	}
}

func TestPasswordHash(t *testing.T) {
	h1 := passwordHash("12345")
	h2 := passwordHash("12345")
	h3 := passwordHash("54321")

	if h1 == "" {
		t.Fatal("hash is empty")
	}

	if h1 != h2 {
		t.Fatal("same password must produce same hash")
	}

	if h1 == h3 {
		t.Fatal("different passwords must produce different hashes")
	}
}

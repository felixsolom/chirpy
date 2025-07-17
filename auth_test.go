package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/felixsolom/chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := auth.HashPassword(password1)
	hash2, _ := auth.HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := auth.CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	uid := uuid.MustParse("6d7b7e37-407f-4002-8d92-67f06a898300")
	secret := "super-secret"

	type testCase struct {
		name       string
		ttl        time.Duration
		checkKey   string
		wantErr    bool
		wantUserID uuid.UUID
	}

	cases := []testCase{
		{
			name:       "happy-path",
			ttl:        time.Hour,
			checkKey:   secret,
			wantErr:    false,
			wantUserID: uid,
		},
		{
			name:     "expired token",
			ttl:      -time.Minute,
			checkKey: secret,
			wantErr:  true,
		},
		{
			name:     "bad secret",
			ttl:      time.Hour,
			checkKey: "other-secret",
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tok, err := auth.MakeJWT(uid, secret, tc.ttl)
			if err != nil {
				t.Fatalf("MakeJWT returned error: %v", err)
			}

			gotID, err := auth.ValidateJWT(tok, tc.checkKey)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("ValidateJWT() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("ValidateJWT() unexpected error: %v", err)
			}
			if gotID != tc.wantUserID {
				t.Fatalf("ValidateJWT() userID mismatch: want %s got %s", tc.wantUserID, gotID)
			}
		})
	}
}

func TestGetBearer(t *testing.T) {

	makeHdr := func(v string) http.Header {
		h := http.Header{}
		if v != "" {
			h.Set("Authorization", v)
		}
		return h
	}

	type TestCase struct {
		name    string
		header  http.Header
		wantErr bool
	}

	cases := []TestCase{
		{
			name:    "valid header",
			header:  makeHdr("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
			wantErr: false,
		},
		{
			name:    "also valid header",
			header:  makeHdr("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"),
			wantErr: false,
		},
		{
			name:    "invalid header",
			header:  makeHdr(""),
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := auth.GetBearerToken(tc.header)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetBearerToken() err = %v, wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

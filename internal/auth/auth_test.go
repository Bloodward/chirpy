package auth

import (
	"net/http"
	"testing"
)

func TestHashAndCheck(t *testing.T) {
	hash, err := HashPassword("04234")
	if err != nil {
		t.Fatalf("HashPassword error : %v", err)
	}

	match, err := CheckPasswordHash("04234", hash)
	if err != nil || !match {
		t.Errorf("expected match, got match=%v err %v", match, err)
	}
	match, err = CheckPasswordHash("wrong", hash)
	if err != nil || match {
		t.Errorf("expected no match, got match=%v err %v", match, err)
	}
}

func TestHashIsSalted(t *testing.T) {
	h1, _ := HashPassword("same")
	h2, _ := HashPassword("same")
	if h1 == h2 {
		t.Errorf("expected different hashes for same password (random salt)")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{"valid", "Bearer abc.def.ghi", "abc.def.ghi", false},
		{"extra whitespace", "Bearer    abc.def.ghi  ", "abc.def.ghi", false},
		{"missing header", "", "", true},
		{"wrong scheme", "Basic abc123", "", true},
		{"bearer no token", "Bearer ", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			if tt.header != "" {
				h.Set("Authorization", tt.header)
			}
			got, err := GetBearerToken(h)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

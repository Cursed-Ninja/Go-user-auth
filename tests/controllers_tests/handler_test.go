package controllers_tests

import (
	"net/http"
	"testing"
	"user-auth/internal/controllers"
)

type mockResponseWriter struct {
	header http.Header
}

func (m *mockResponseWriter) Header() http.Header {
	if m.header == nil {
		m.header = http.Header{}
	}
	return m.header
}

func (m *mockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
}

// This is a dummy test and will fail as SQL database is not connected
// Also form type is not mulitpart/form-data 
func TestLoginUserHandler(t *testing.T) {
	tests := []struct {
		name string
		form map[string][]string
	}{
		{
			name: "Valid login",
			form: map[string][]string{
				"email":    {"john.doe@email.com"},
				"password": {"password"},
			},
		},
		{
			name: "User does not exists",
			form: map[string][]string{
				"email":    {"abc@email.com"},
				"password": {"password"},
			},
		},
		{
			name: "Invalid password",
			form: map[string][]string{
				"email":    {"john.doe@email.com"},
				"password": {"password123"},
			},
		},
	}

	for _, test := range tests {
		req, err := http.NewRequest("POST", "/login", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Form = test.form
		mockWriter := &mockResponseWriter{}

		controllers.LoginUserHandler(mockWriter, req)

		if mockWriter.Header().Get("Location") != "/profile" {
			t.Error("Expected redirect to /profile")
		}
	}
}

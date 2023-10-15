package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"errors"

	"demerzel-badges/middleware"
)

type MockExternalService struct {
	mock.Mock
}

func (m *MockExternalService) GetUserID(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func TestAuthMiddleware(t *testing.T) {
	r := gin.Default()
	externalService := &MockExternalService{}

	r.GET("/protected", middleware.AuthMiddleware(externalService), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	tests := []struct {
		Name            string
		Authorization   string
		ExpectedStatus  int
		ExpectedMessage string
	}{
		{"Valid Token", "Bearer valid-token", http.StatusOK, "Success"},
		{"Invalid Token", "Bearer invalid-token", http.StatusUnauthorized, "Invalid or expired token"},
		{"Missing Authorization", "", http.StatusUnauthorized, "Missing Authorization header"},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", test.Authorization)

			token := strings.TrimPrefix(test.Authorization, "Bearer ")

			if test.ExpectedStatus == http.StatusOK {
				externalService.On("GetUserID", token).Return("user_id", nil)
			} else {
				externalService.On("GetUserID", token).Return("", errors.New("Invalid token"))
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != test.ExpectedStatus {
				t.Errorf("Expected status %d, got %d", test.ExpectedStatus, w.Code)
			}

			expectedResponse := `{"message":"` + test.ExpectedMessage + `"}`
			if w.Body.String() != expectedResponse {
				t.Errorf("Expected response %s, got %s", expectedResponse, w.Body.String())
			}
		})
	}
}

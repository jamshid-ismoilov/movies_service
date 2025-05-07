package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"movies_service/model"
	"movies_service/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// stubUserService is a stub implementation of UserService for handler tests
type stubUserService struct {
	LoginFn    func(username, password string) (string, error)
	RegisterFn func(username, password string) (*model.User, error)
}

func (s *stubUserService) Login(username, password string) (string, error) {
	return s.LoginFn(username, password)
}
func (s *stubUserService) Register(username, password string) (*model.User, error) {
	return s.RegisterFn(username, password)
}

func TestUserHandler_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stubService := &stubUserService{
		LoginFn: func(u, p string) (string, error) {
			return "", service.ErrInvalidCredentials
		},
		RegisterFn: func(u, p string) (*model.User, error) {
			return nil, nil
		},
	}
	handler := NewUserHandler(stubService)

	body := []byte(`{"username":"alice","password":"wrong"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Login(c)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "invalid username or password", resp["error"])
}

func TestUserHandler_Register_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	stubService := &stubUserService{
		RegisterFn: func(u, p string) (*model.User, error) {
			return &model.User{ID: 1, Username: u, Password: ""}, nil
		},
		LoginFn: func(u, p string) (string, error) {
			return "", nil
		},
	}
	handler := NewUserHandler(stubService)

	body := []byte(`{"username":"bob","password":"secret"}`)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.Register(c)

	require.Equal(t, http.StatusCreated, w.Code)
	var user model.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	require.NoError(t, err)
	require.Equal(t, "bob", user.Username)
	require.Equal(t, uint(1), user.ID)
	require.Empty(t, user.Password)
}

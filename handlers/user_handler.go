package handlers

import (
	"net/http"

	"movies_service/model"
	"movies_service/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body model.User true "User credentials"
// @Success 201 {object} model.User
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}
	created, err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		if err == service.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "username already taken"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		}
		return
	}
	c.JSON(http.StatusCreated, created)
}

// Login godoc
// @Summary Log in a user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body model.User true "User credentials"
// @Success 200 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data"})
		return
	}
	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not login"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

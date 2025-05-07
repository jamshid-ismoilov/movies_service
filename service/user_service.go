package service

import (
	"errors"

	"movies_service/auth"
	"movies_service/model"
	"movies_service/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// defining service-level errors
var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotFound           = errors.New("not found")
)

type UserService interface {
	Register(username, password string) (*model.User, error)
	Login(username, password string) (string, error)
}

type userServiceImpl struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewUserService(userRepo repository.UserRepository, jwtSecret string) UserService {
	return &userServiceImpl{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *userServiceImpl) Register(username, password string) (*model.User, error) {
	// checking if user already exists
	if _, err := s.userRepo.GetByUsername(username); err == nil {
		return nil, ErrUserExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username: username,
		Password: string(hashed),
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *userServiceImpl) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}
	token, err := auth.GenerateToken(user, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

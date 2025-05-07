package service

import (
	"testing"

	"movies_service/model"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// fakeUserRepo is a fake implementation of UserRepository for tests
type fakeUserRepo struct {
	users  map[string]model.User
	lastID uint
}

func (f *fakeUserRepo) Create(user *model.User) error {
	f.lastID++
	user.ID = f.lastID
	f.users[user.Username] = *user
	return nil
}

func (f *fakeUserRepo) GetByUsername(username string) (*model.User, error) {
	u, ok := f.users[username]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	userCopy := u
	return &userCopy, nil
}

func (f *fakeUserRepo) GetByID(id uint) (*model.User, error) {
	for _, u := range f.users {
		if u.ID == id {
			userCopy := u
			return &userCopy, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: make(map[string]model.User)}
}

func TestUserService_RegisterAndLogin(t *testing.T) {
	repo := newFakeUserRepo()
	secret := "testsecret"
	service := NewUserService(repo, secret)

	// Register a new user
	user, err := service.Register("jamshid", "password123")
	require.NoError(t, err, "register should succeed for new user")
	require.NotNil(t, user)
	require.Equal(t, "jamshid", user.Username)
	require.NotZero(t, user.ID)
	require.Equal(t, "", user.Password)

	_, err = service.Register("jamshid", "newpass")
	require.Error(t, err)
	require.Equal(t, ErrUserExists, err, "should error that user exists")

	token, err := service.Login("jamshid", "password123")
	require.NoError(t, err, "login with correct password should succeed")
	require.NotEmpty(t, token, "token should be returned")

	_, err = service.Login("jamshid", "wrongpass")
	require.Error(t, err)
	require.Equal(t, ErrInvalidCredentials, err, "should get invalid credentials error")
}

func TestUserService_PasswordHashing(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo, "secret")
	username := "bob"
	rawPassword := "mypassword"
	user, err := svc.Register(username, rawPassword)
	require.NoError(t, err)
	require.NotNil(t, user)
	stored, err := repo.GetByUsername(username)
	require.NoError(t, err)
	require.NotEqual(t, rawPassword, stored.Password, "stored password should be hashed")
	err = bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(rawPassword))
	require.NoError(t, err, "stored password hash should match original password")
}

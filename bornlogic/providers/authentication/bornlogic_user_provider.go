package authentication

import (
	"context"
	"github.com/authelia/authelia/v4/bornlogic/repositories"
	"github.com/authelia/authelia/v4/internal/authentication"
)

type BornlogicUserProvider struct {
	userRepository repositories.User
}

func NewBornlogicUserProvider(userRepository repositories.User) authentication.UserProvider {
	return &BornlogicUserProvider{userRepository: userRepository}
}

// StartupCheck implements the startup check provider interface.
func (b BornlogicUserProvider) StartupCheck() (err error) {
	return nil
}

// CheckUserPassword checks if provided password matches for the given user.
func (b BornlogicUserProvider) CheckUserPassword(username string, password string) (valid bool, err error) {
	user, err := b.userRepository.Get(context.Background(), username)
	if err != nil {
		return false, handleBusinessError(err)
	}
	return CheckPassword(user.HashedPassword, password)
}

func handleBusinessError(err error) error {
	switch err {
	case repositories.ErrNotFound:
		return authentication.ErrUserNotFound
	}
	return err
}

// GetDetails retrieve the groups a user belongs to.
func (b BornlogicUserProvider) GetDetails(username string) (details *authentication.UserDetails, err error) {
	user, err := b.userRepository.Get(context.Background(), username)
	if err != nil {
		return nil, handleBusinessError(err)
	}
	return &authentication.UserDetails{
		Username:    username,
		DisplayName: user.Name,
		Emails:      []string{user.Email},
		Groups:      []string{"default"},
	}, nil
}

// UpdatePassword update the password of the given user.
func (b BornlogicUserProvider) UpdatePassword(username string, newPassword string) (err error) {
	panic("implement me")
}

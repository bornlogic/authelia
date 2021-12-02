package postgres

import (
	"context"
	"github.com/authelia/authelia/v4/bornlogic/entities"
	"github.com/authelia/authelia/v4/bornlogic/infrastructure/database/postgres"
	"github.com/authelia/authelia/v4/bornlogic/repositories"
	"gorm.io/gorm"
)

type userPostgresRepository struct {
	*postgres.Client
}

func NewUserRepository(client *postgres.Client) repositories.User {
	return &userPostgresRepository{
		Client: client,
	}
}

func (u userPostgresRepository) Get(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	if err := u.Gorm.Where("USERNAME = ?", username).First(&user).Error; err != nil {
		return nil, handleError(err)
	}
	return &user, nil
}

func handleError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return repositories.ErrNotFound
	}
	return err
}

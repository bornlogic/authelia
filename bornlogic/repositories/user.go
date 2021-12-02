package repositories

import (
	"context"
	"errors"
	"github.com/authelia/authelia/v4/bornlogic/entities"
)

var ErrNotFound = errors.New("user not found")

type User interface {
	Get(ctx context.Context, username string) (*entities.User, error)
}

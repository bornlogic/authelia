package providers

import (
	"github.com/authelia/authelia/v4/bornlogic/configuration/environment"
	"github.com/authelia/authelia/v4/bornlogic/entities"
	"github.com/authelia/authelia/v4/bornlogic/infrastructure/database"
	bornlogicAuth "github.com/authelia/authelia/v4/bornlogic/providers/authentication"
	repositories "github.com/authelia/authelia/v4/bornlogic/repositories/postgres"
	"github.com/authelia/authelia/v4/internal/authentication"
	"github.com/sirupsen/logrus"
)

type Provider struct {
	UserProvider authentication.UserProvider
}

func NewProvider() (*Provider, error) {
	log := logrus.New()
	log.Infoln("Starting Bornlogic provider")

	configuration, err := environment.Load()
	if err != nil {
		log.Fatalln("Error loading environment", err)
	}

	pgClient, err := database.NewPostgreSQLClient(configuration.PostgreSQL)
	if err != nil {
		log.Fatalln("Error loading environment", err)
	}
	if err := pgClient.AutoMigrate(&entities.User{}); err != nil {
		log.Fatalln("Error auto migrating database", err)
	}

	//providers
	userProvider := bornlogicAuth.NewBornlogicUserProvider(
		repositories.NewUserRepository(pgClient))

	return &Provider{
		UserProvider: userProvider,
	}, nil
}

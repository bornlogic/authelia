package bornlogic

import (
	"github.com/authelia/authelia/v4/bornlogic/configuration/environment"
	"github.com/authelia/authelia/v4/bornlogic/entities"
	"github.com/authelia/authelia/v4/bornlogic/infrastructure/database/postgres"
	repositories "github.com/authelia/authelia/v4/bornlogic/repositories/postgres"
	bornlogicAuth "github.com/authelia/authelia/v4/bornlogic/services/authentication"
	"github.com/authelia/authelia/v4/internal/authentication"
	"github.com/sirupsen/logrus"
)

type Provider struct {
	UserProvider authentication.UserProvider
}

func NewProvider() (*Provider, error) {
	log := logrus.New()
	log.Infoln("Starting bornlogic provider")

	configuration, err := environment.Load()
	if err != nil {
		log.Fatalln("Error loading environment", err)
	}

	postgresClient, err := postgres.NewClient(configuration.Postgres)
	if err != nil {
		log.Fatalln("Error loading environment", err)
	}
	if err := postgresClient.AutoMigrate(&entities.User{}); err != nil {
		log.Fatalln("Error auto migrating database", err)
	}

	//repositories
	userRepository := repositories.NewUserRepository(postgresClient)

	//providers
	userProvider := bornlogicAuth.NewBornlogicUserProvider(userRepository)

	return &Provider{
		UserProvider: userProvider,
	}, nil
}

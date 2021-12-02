package commands

import (
	bornlogicProvider "github.com/authelia/authelia/v4/bornlogic/providers"
	"github.com/authelia/authelia/v4/internal/authentication"
	"github.com/authelia/authelia/v4/internal/authorization"
	"github.com/authelia/authelia/v4/internal/middlewares"
	"github.com/authelia/authelia/v4/internal/notification"
	"github.com/authelia/authelia/v4/internal/ntp"
	"github.com/authelia/authelia/v4/internal/oidc"
	"github.com/authelia/authelia/v4/internal/regulation"
	"github.com/authelia/authelia/v4/internal/session"
	"github.com/authelia/authelia/v4/internal/storage"
	"github.com/authelia/authelia/v4/internal/utils"
)

func getStorageProvider() (provider storage.Provider) {
	switch {
	case config.Storage.PostgreSQL != nil:
		return storage.NewPostgreSQLProvider(*config.Storage.PostgreSQL, config.Storage.EncryptionKey)
	case config.Storage.MySQL != nil:
		return storage.NewMySQLProvider(*config.Storage.MySQL, config.Storage.EncryptionKey)
	case config.Storage.Local != nil:
		return storage.NewSQLiteProvider(config.Storage.Local.Path, config.Storage.EncryptionKey)
	default:
		return nil
	}
}

func getProviders() (providers middlewares.Providers, warnings []error, errors []error) {
	// TODO: Adjust this so the CertPool can be used like a provider.
	autheliaCertPool, warnings, errors := utils.NewX509CertPool(config.CertificatesDirectory)
	if len(warnings) != 0 || len(errors) != 0 {
		return providers, warnings, errors
	}

	storageProvider := getStorageProvider()

	var (
		userProvider authentication.UserProvider
		err          error
	)

	switch {
	case config.AuthenticationBackend.File != nil:
		userProvider = authentication.NewFileUserProvider(config.AuthenticationBackend.File)
	case config.AuthenticationBackend.LDAP != nil:
		userProvider = authentication.NewLDAPUserProvider(config.AuthenticationBackend, autheliaCertPool)
	}

	//TODO fix me
	bornlogicProvider, err := bornlogicProvider.NewProvider()
	if err != nil {
		panic(err)
	}
	userProvider = bornlogicProvider.UserProvider

	var notifier notification.Notifier

	switch {
	case config.Notifier.SMTP != nil:
		notifier = notification.NewSMTPNotifier(config.Notifier.SMTP, autheliaCertPool)
	case config.Notifier.FileSystem != nil:
		notifier = notification.NewFileNotifier(*config.Notifier.FileSystem)
	}

	var ntpProvider *ntp.Provider
	if config.NTP != nil {
		ntpProvider = ntp.NewProvider(config.NTP)
	}

	clock := utils.RealClock{}
	authorizer := authorization.NewAuthorizer(config)
	sessionProvider := session.NewProvider(config.Session, autheliaCertPool)
	regulator := regulation.NewRegulator(config.Regulation, storageProvider, clock)

	oidcProvider, err := oidc.NewOpenIDConnectProvider(config.IdentityProviders.OIDC)
	if err != nil {
		errors = append(errors, err)
	}

	return middlewares.Providers{
		Authorizer:      authorizer,
		UserProvider:    userProvider,
		Regulator:       regulator,
		OpenIDConnect:   oidcProvider,
		StorageProvider: storageProvider,
		NTP:             ntpProvider,
		Notifier:        notifier,
		SessionProvider: sessionProvider,
	}, warnings, errors
}

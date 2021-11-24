package suites

import (
	"fmt"
	"os"
	"time"
)

var cliSuiteName = "CLI"

func init() {
	_ = os.MkdirAll("/tmp/authelia/CLISuite/", 0700)

	dockerEnvironment := NewDockerEnvironment([]string{
		"internal/suites/docker-compose.yml",
		"internal/suites/CLI/docker-compose.yml",
		"internal/suites/example/compose/authelia/docker-compose.backend.{}.yml",
	})

	setup := func(suitePath string) error {
		if err := dockerEnvironment.Up(); err != nil {
			return err
		}

		return waitUntilAutheliaIsReady(dockerEnvironment, cliSuiteName)
	}

	displayAutheliaLogs := func() error {
		backendLogs, err := dockerEnvironment.Logs("authelia-backend", nil)
		if err != nil {
			return err
		}

		fmt.Println(backendLogs)

		return nil
	}

	teardown := func(suitePath string) error {
		err := dockerEnvironment.Down()
		_ = os.Remove("/tmp/db.suites.sqlite3")

		return err
	}

	GlobalRegistry.Register(cliSuiteName, Suite{
		SetUp:           setup,
		SetUpTimeout:    5 * time.Minute,
		OnSetupTimeout:  displayAutheliaLogs,
		TestTimeout:     2 * time.Minute,
		TearDown:        teardown,
		TearDownTimeout: 2 * time.Minute,
	})
}

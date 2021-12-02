package postgres

import (
	"database/sql"
	"github.com/authelia/authelia/v4/bornlogic/configuration/environment"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const driverName = "postgres"

type Client struct {
	*sql.DB
	Gorm *gorm.DB
}

func NewClient(config environment.Postgres) (*Client, error) {
	db, err := sql.Open(driverName, config.BuildConnectionString())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Client{
		DB:   db,
		Gorm: gormDB,
	}, nil
}

func (c *Client) AutoMigrate(entities ...interface{}) error {
	return c.Gorm.AutoMigrate(entities...)
}

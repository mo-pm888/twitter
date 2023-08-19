package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	ServerHost string
	ServerPort string
	DbHost     string
}

func New() (*Config, error) {
	var c = new(Config)
	err := c.LoadEnvFile()
	return c, err
}

func (c *Config) LoadEnvFile() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	c.DbPort = os.Getenv("DB_PG_PORT")
	c.DbUser = os.Getenv("DB_PG_USER")
	c.DbPassword = os.Getenv("DB_PG_PASSWORD")
	c.DbName = os.Getenv("DB_PG_NAME")
	c.DbHost = os.Getenv("SERVER_HOST")
	c.ServerPort = os.Getenv("SERVER_PORT")
	c.ServerHost = os.Getenv("DB_PG_HOST")
	return nil
}

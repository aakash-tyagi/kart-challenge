package config

// use dot env to load the config file
import (
	"os"

	"github.com/joho/godotenv"
)

// Config is a struct that holds the configuration for the application
type Config struct {
	DBUrl        string
	ServerPort   string
	DatabaseName string
}

// LoadConfig loads the config file into the Config struct
func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	config := Config{
		DBUrl:        os.Getenv("DB_CONN_STRING"),
		ServerPort:   os.Getenv("SERVER_PORT"),
		DatabaseName: os.Getenv("DB_NAME"),
	}

	return &config, nil
}

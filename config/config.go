package config


import (
		"os"
		_ "github.com/joho/godotenv/autoload"
		"fmt"
	)	

type Config struct {
}

func (c *Config) GetConfig() map[string]string {

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	db_config := map[string]string{
		"dev": fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD")),
	}

	// should be moved to environment variables once the main checklist is done.
	config := map[string]string{
		"foo":        "Bar",
		"jwtKey":     "SomeRandomJWTKEY",
		"db_dialect": "postgres",
		"db_env":     db_config["dev"],
		"PORT":       port,
	}
	return config
}

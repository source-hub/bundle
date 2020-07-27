package config


import "os"

type Config struct {
}

func (c *Config) GetConfig() map[string]string {

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}

	db_env := map[string]string{
		"aakash":   "host=postgres_12_container port=5432 user=konem dbname=bundle password=PassWord sslmode=disable",
	}

	// should be moved to environment variables once the main checklist is done.
	config := map[string]string{
		"foo":        "Bar",
		"jwtKey":     "SomeRandomJWTKEY",
		"db_dialect": "postgres",
		"db_env":     db_env["aakash"],
		"sendgrid":   "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",
		"PORT":       port,
	}
	return config
}

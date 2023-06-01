package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port				string
	DBHost			string
	DBPort			string
	DBName			string
	DBUser     	string
	DBPassword	string
	// Add more configuration variables as needed
}

var Env = Config{}

// load env configuration
func InitConfig() {
	// Load the .env file
	goEnv := os.Getenv("GO_ENV");
	
	if (goEnv != "production") {
		// Load env from `.env` file 
		err := godotenv.Load()
		if err != nil {
				log.Fatal("Error loading .env file")
		}
	}

	// Create a Config instance and populate it with the environment variables
	Env.Port = os.Getenv("PORT");
	if (Env.Port == "") {
		// Set default port
		Env.Port = "8080"
	}
	Env.DBHost = os.Getenv("DB_HOST");
	Env.DBPort = os.Getenv("DB_PORT");
	Env.DBName = os.Getenv("DB_NAME");
	Env.DBUser = os.Getenv("DB_USER");
	Env.DBPassword = os.Getenv("DB_PASSWORD");
}

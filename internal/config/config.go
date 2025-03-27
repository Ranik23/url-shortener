package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
)


type DBConfig struct {
	Type		string	`yaml:"type"`
	Host 		string	`yaml:"host"`
	Port 		string	`yaml:"port"`
	UserName 	string	`yaml:"username"`
	Password 	string	
	DBName 		string	
	SSL 		string	`yaml:"ssl"`
}


type HTTPServerConfig struct {
	Host		string	`yaml:"host"`
	Port		string	`yaml:"port"`
}


type GRPCServerConfig struct {
	Host		string	`yaml:"host"`
	Port		string	`yaml:"port"`
}


type Config struct {
	Database 	 DBConfig			`yaml:"database"`
	HTTPServer   HTTPServerConfig	`yaml:"http_server"`
	GRPCServer	 GRPCServerConfig	`yaml:"grpc_server"`
}

func LoadConfig(envPath string, configPath string) (*Config, error) {

	if err := godotenv.Load(envPath); err != nil {
		return nil, err
	}

	file, err := os.ReadFile(configPath)
    if err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }

	var config Config
    err = yaml.Unmarshal(file, &config)
    if err != nil {
        log.Fatalf("Error unmarshalling YAML: %v", err)
    }

	config.Database.UserName = os.Getenv("DB_USER_NAME")
    config.Database.Password = os.Getenv("DB_PASSWORD")

	return &config, nil
}
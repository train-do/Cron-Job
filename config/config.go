package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	AppDebug    bool
	AppSecret   string
	ServerPort  string
	DBMigrate   bool
	DBSeeding   bool
	RedisConfig RedisConfig
}

type RedisConfig struct {
	Url      string
	Password string
	Prefix   string
}

func LoadConfig(migrateDb bool, seedDb bool) (Config, error) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.SetConfigType("dotenv")
	viper.SetConfigName(".env")

	// Set default values
	setDefaultValues(migrateDb, seedDb)

	// Allow Viper to read environment variables
	viper.AutomaticEnv()

	// Read the configuration file
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file: %s, using default values or environment variables", err)
	}

	// add value to the config
	config := Config{
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		AppDebug:   viper.GetBool("APP_DEBUG"),
		AppSecret:  viper.GetString("APP_SECRET"),
		ServerPort: viper.GetString("SERVER_PORT"),
		DBMigrate:  viper.GetBool("DB_MIGRATE"),
		DBSeeding:  viper.GetBool("DB_SEEDING"),
		RedisConfig: RedisConfig{
			Url:      viper.GetString("REDIS_URL"),
			Password: viper.GetString("REDIS_PASSWORD"),
			Prefix:   viper.GetString("REDIS_PREFIX"),
		},
	}
	return config, nil
}

func setDefaultValues(migrateDb bool, seedDb bool) {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "admin")
	viper.SetDefault("DB_NAME", "database")
	viper.SetDefault("APP_DEBUG", true)
	viper.SetDefault("APP_SECRET", "team-1")
	viper.SetDefault("SERVER_PORT", ":8080")

	viper.SetDefault("DB_MIGRATE", migrateDb)
	viper.SetDefault("DB_SEEDING", seedDb)
}

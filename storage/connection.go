package storage

import (
	"backtest/logger"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	host = "postgres" //change to localhost if running locally
	port = "5432"
)

func GetFromEnv(env_name string) string {
	meta, exists := os.LookupEnv(env_name)
	if !exists {
		logger.Err.Fatalln("No such env variable found")
	}
	return meta
}

func GetDBConnection() *sql.DB {
	if err := godotenv.Load(); err != nil {
		logger.Err.Fatalln("No .env file found")
	}

	db_username := GetFromEnv("DB_USERNAME")
	db_password := GetFromEnv("DB_PASSWORD")
	db_name := GetFromEnv("DB_NAME")

	connection_info := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, db_username, db_password, db_name)

	db, err := sql.Open("postgres", connection_info)

	if err != nil {
		logger.Err.Fatalln("Next error occured during DB connection:", err)
	}

	logger.Info.Printf("Established connection with PostreSQL database - Name: %s, port: %s", db_name, port)
	return db
}

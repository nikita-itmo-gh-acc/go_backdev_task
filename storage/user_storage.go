package storage

import (
	"backtest/logger"
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	Id       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
}

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Get(id uuid.UUID) (User, error) {
	var user User
	err := s.db.QueryRow("SELECT * FROM Users WHERE id = $1", id).Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		logger.Err.Fatalln("No such user found")
	}
	return user, err
}

func (s *UserStorage) Create(model *User) error {
	_, err := s.db.Exec("INSERT INTO Users (id, email, password) VALUES ($1, $2, $3)", model.Id, model.Email, model.Password)
	if err != nil {
		logger.Err.Fatalln("Can't create user")
	}
	logger.Info.Println("CREATED user with ID =", model.Id)
	return err
}

func (s *UserStorage) Delete(model *User) error {
	_, err := s.db.Exec(`DELETE FROM Users WHERE id = $1`, model.Id)
	if err != nil {
		logger.Err.Fatalln("Can't delete user", err)
	}
	logger.Info.Println("DELETED user with ID =", model.Id)
	return err
}

func CreateUserStorage(conn *sql.DB) *UserStorage {
	_, err := conn.Exec(`CREATE TABLE IF NOT EXISTS Users (
						id uuid PRIMARY KEY,
						email varchar(255) UNIQUE NOT NULL,
						password varchar(255) NOT NULL
	)`)

	if err != nil {
		logger.Err.Fatalln("Next error occured during Users Table creation:", err)
	}

	return &UserStorage{
		db: conn,
	}
}

package storage

import (
	"backtest/logger"
	"database/sql"

	"github.com/google/uuid"
)

type Session struct {
	Id           uuid.UUID `json:"session_id,omitempty"`
	UserId       uuid.UUID `json:"user_id"`
	RefreshToken []byte    `json:"-"`
	AccessToken  string    `json:"access_token,omitempty"`
	Ip           string    `json:"ip_address"`
	// ExpiresAt    time.Time `JSON:"expires_at"`
}

type SessionStorage struct {
	db *sql.DB
}

func (s *SessionStorage) Get(id uuid.UUID) (Session, error) {
	var session Session
	err := s.db.QueryRow(`SELECT id, userId, refreshToken, ip
        FROM Sessions WHERE id = $1`, id).Scan(&session.Id, &session.UserId, &session.RefreshToken, &session.Ip)
	if err != nil {
		logger.Err.Fatalln("No such session found")
	}
	return session, err
}

func (s *SessionStorage) Create(model *Session) error {
	_, err := s.db.Exec(`INSERT INTO Sessions (id, userId, refreshToken, ip)
                         VALUES ($1, $2, $3, $4)`, model.Id, model.UserId, model.RefreshToken, model.Ip)
	if err != nil {
		logger.Err.Fatalln("Can't create session", err)
	}
	logger.Info.Println("CREATED session with ID =", model.Id)
	return err
}

func (s *SessionStorage) Delete(model *Session) error {
	_, err := s.db.Exec(`DELETE FROM Sessions WHERE id = $1`, model.Id)
	if err != nil {
		logger.Err.Fatalln("Can't delete session", err)
	}
	logger.Info.Println("DELETED session with ID =", model.Id)
	return err
}

func CreateSessionStorage(conn *sql.DB) *SessionStorage {
	_, err := conn.Exec(`CREATE OR REPLACE FUNCTION get_future_date(days INT)
                        RETURNS timestamptz
                        LANGUAGE plpgsql
                        AS $$
                        BEGIN
                            RETURN NOW() + INTERVAL '1 day' * days;
                        END;
                        $$;`)
	if err != nil {
		logger.Err.Fatalln("Next error occured during get_future_date() creation:", err)
	}

	_, err = conn.Exec(`CREATE TABLE IF NOT EXISTS Sessions (
                        id uuid PRIMARY KEY,
                        userId uuid REFERENCES Users(id) ON DELETE CASCADE,
                        refreshToken BYTEA NOT NULL,
                        ip varchar(15) NOT NULL,
                        createdAt timestamp with time zone NOT NULL DEFAULT get_future_date(0),
                        expiresAt timestamp NOT NULL DEFAULT get_future_date(30)
    )`)

	if err != nil {
		logger.Err.Fatalln("Next error occured during Sessions Table creation:", err)
	}

	return &SessionStorage{
		db: conn,
	}
}

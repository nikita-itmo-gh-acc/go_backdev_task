package main

import (
	"backtest/logger"
	"backtest/storage"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"

	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	port      string = ":4242"
	smtp_host string = "smtp.gmail.com"
	smtp_port string = ":587"
)

func init() {
	logger.InitLoggers()
}

func main() {
	db_conn := storage.GetDBConnection()
	defer db_conn.Close()

	home_handler := NewHomeHandler()
	auth_handler := NewAuthHandler(db_conn)

	mux := http.NewServeMux()
	mux.Handle("/", home_handler)
	mux.Handle("/protected", home_handler)
	mux.Handle("/api/auth/generate-tokens", auth_handler)
	mux.Handle("/api/auth/refresh-tokens", auth_handler)

	logger.Info.Printf("Запуск сервера на %s", storage.GetFromEnv("HOST")+port)
	http.ListenAndServe(port, mux)
}

type AccessJWTClaims struct {
	IpAdress string
	jwt.RegisteredClaims
}

func GenerateAccessToken(user_id uuid.UUID, ip_address string) (string, error) {
	var secret_key = []byte(storage.GetFromEnv("JWT_SECRET_KEY"))
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, AccessJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user_id.String(),
			Issuer:    "backtest",
			Audience:  []string{"user"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30))},
		IpAdress: ip_address,
	})

	access_token_str, err := claims.SignedString(secret_key)
	if err != nil {
		logger.Err.Fatalln("Error occured during access token creation:", err)
	}
	return access_token_str, err
}

func VerifyAccessToken(token_string string) error {
	var secret_key = []byte(storage.GetFromEnv("JWT_SECRET_KEY"))
	token, err := jwt.ParseWithClaims(token_string, &AccessJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret_key, nil
	})

	if err != nil {
		logger.Err.Println("Can't parse access token", err)
		return err
	}

	if !token.Valid {
		logger.Err.Println("Invalid access token")
		return err
	}
	return nil
}

func PasreJSON(r *http.Request, object interface{}) {
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, object); err != nil {
		logger.Err.Fatalln("Can't decode request body for object creation:", err)
		panic(err)
	}
}

func RenderJSON(w http.ResponseWriter, object interface{}) {
	js, err := json.Marshal(object)
	if err != nil {
		logger.Err.Fatalln("Can't render JSON from object:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func SendEmail(email string, message string) {
	sender := storage.GetFromEnv("MAIL_ADDR")
	pass := storage.GetFromEnv("MAIL_PASS")
	auth := smtp.PlainAuth("", sender, pass, smtp_host)

	err := smtp.SendMail(smtp_host+smtp_port, auth, sender, []string{email}, []byte(message))
	if err != nil {
		logger.Err.Printf("Can't send email to %s\n", email)
	}
	logger.Info.Println("Successfully sent mail to", email)
}

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		w.Write([]byte("Home page!"))
	case r.URL.Path == "/protected":
		h.GetProtectedPage(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (h *HomeHandler) GetProtectedPage(w http.ResponseWriter, r *http.Request) {
	var session storage.Session
	PasreJSON(r, &session)
	if err := VerifyAccessToken(session.AccessToken); err != nil {
		logger.Err.Println("Access token verification failed..")
		return
	}
	json_resp := map[string]string{"message": "You've successfully visited protected page"}
	RenderJSON(w, json_resp)
}

type AuthHandler struct {
	Users    storage.Storage[storage.User]
	Sessions storage.Storage[storage.Session]
}

func NewAuthHandler(conn *sql.DB) *AuthHandler {
	user_storage := storage.CreateUserStorage(conn)
	session_storage := storage.CreateSessionStorage(conn)

	return &AuthHandler{
		Users:    user_storage,
		Sessions: session_storage,
	}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		w.Write([]byte("got!"))
	case r.Method == http.MethodPost && r.URL.Path == "/api/auth/generate-tokens":
		h.CreateSession(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/api/auth/refresh-tokens":
		h.RefreshSession(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *AuthHandler) CreateSession(w http.ResponseWriter, r *http.Request) {
	var (
		session           storage.Session
		session_id        uuid.UUID = uuid.New()
		new_refresh_token uuid.UUID = uuid.New()
	)

	logger.Info.Println("Session initialization process begin...")
	PasreJSON(r, &session)
	fmt.Println("User ID:", session.UserId)
	session.Id = session_id
	hashed_refresh_token, err := bcrypt.GenerateFromPassword([]byte(new_refresh_token.String()), bcrypt.DefaultCost)

	if err != nil {
		logger.Err.Println("Error occured while hashing refresh token:", err)
		return
	}

	session.RefreshToken = hashed_refresh_token
	h.Sessions.Create(&session)

	logger.Info.Println("Session successfully initialized!")

	access_token, _ := GenerateAccessToken(session.UserId, session.Ip)
	session.AccessToken = access_token

	cookie := http.Cookie{
		Name:     "RefreshTokenCookie",
		Value:    base64.StdEncoding.EncodeToString([]byte(new_refresh_token.String())),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("Status", "Created")
	RenderJSON(w, session)
}

func (h *AuthHandler) RefreshSession(w http.ResponseWriter, r *http.Request) {
	var (
		new_session    storage.Session
		new_session_id uuid.UUID = uuid.New()
	)
	logger.Info.Println("Refresh tokens process begin...")

	PasreJSON(r, &new_session)
	prev_session, _ := h.Sessions.Get(new_session.Id)

	refresh_token_cookie, err := r.Cookie("RefreshTokenCookie")
	if err != nil {
		logger.Err.Println("Can't get refresh token from request:", err)
		return
	}

	refresh_token, err := base64.StdEncoding.DecodeString(refresh_token_cookie.Value)
	if err != nil {
		logger.Err.Println("Can't get decode refresh token:", err)
		return
	}

	err = bcrypt.CompareHashAndPassword(prev_session.RefreshToken, refresh_token)
	if err != nil {
		logger.Err.Println("refresh tokens aren't equal")
		return
	}

	new_refresh_token := uuid.New()
	new_hashed_refresh_token, err := bcrypt.GenerateFromPassword([]byte(new_refresh_token.String()), bcrypt.DefaultCost)
	if err != nil {
		logger.Err.Println("Error occured during refresh token hashing:", err)
		return
	}
	new_session.RefreshToken = new_hashed_refresh_token
	new_session.Id = new_session_id
	new_session.UserId = prev_session.UserId

	access_token, _ := GenerateAccessToken(new_session.UserId, new_session.Ip)
	new_session.AccessToken = access_token

	if prev_session.Ip != new_session.Ip {
		user, err := h.Users.Get(new_session.UserId)
		if err != nil {
			logger.Err.Println("Can't get user:", err)
			return
		}
		SendEmail(user.Email, fmt.Sprintf("Someone trying to get access to your account from %s", new_session.Ip))
	}

	h.Sessions.Delete(&prev_session)
	h.Sessions.Create(&new_session)

	logger.Info.Println("Tokens successfully updated!")

	cookie := http.Cookie{
		Name:     "RefreshTokenCookie",
		Value:    base64.StdEncoding.EncodeToString([]byte(new_refresh_token.String())),
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("Status", "OK")
	RenderJSON(w, new_session)
}

package data

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func init() {
	var sqlite_path string

	switch os.Getenv("ENV") {
	case "LOCAL":
		sqlite_path = "data/local/storage.db"
	case "AZURE":
		sqlite_path = "/home/sqlite/gpteacher/storage.db"
	default:
		log.Fatal("Invalid ENV environment variable; set it as LOCAL or AZURE")
	}

	var err error

	db, err = sql.Open("sqlite3", sqlite_path)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatal("Error pinging database - ", err)
	}
}

// Closes database.
func Close() {
	if db != nil {
		db.Close()
	}
}

func User_Create(username, password_hash string) error {
	var uuid [16]byte

	_, err := rand.Read(uuid[:])
	if err != nil {
		return err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // NOTE: Set the version to UUID v4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // NOTE: Set the variant to RFC 4122

	uuid_string := fmt.Sprintf("%X-%X-%X-%X-%X", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])

	_, err = db.Exec(`
		INSERT INTO Users
			(Id, Username, Password_Hash) 
		VALUES 
			(?, ?, ?)`,
		uuid_string, username, password_hash)

	return err
}

func User_Read() error {
	// TODO
	return errors.New("Function not yet implemented")
}

func User_ReadPasswordHash(username string) (string, error) {
	var password_hash string

	if err := db.QueryRow("SELECT Password_Hash FROM Users WHERE Username = ?", username).Scan(&password_hash); err != nil {
		return "", err
	}

	return password_hash, nil
}

func User_Update() error {
	// TODO
	return errors.New("Function not yet implemented")
}

func User_UpdatePassword(username, new_hash string) error {
	_, err := db.Exec("UPDATE Users SET Password_Hash = ? WHERE Username = ?", new_hash, username)

	return err
}

func User_Delete() error {
	// TODO
	return errors.New("Function not yet implemented")
}

func UserSession_CreateOrUpdate(username, session_token string) error {
	var err error
	var exists int

	if err = db.QueryRow("SELECT 1 FROM UserSessions WHERE Username = ?", username).Scan(&exists); err != nil {
		return err
	}

	if exists == 0 {
		_, err = db.Exec("INSERT INTO UserSessions (Username, SessionToken, Expires_At) VALUES (?, ?, ?)", username, session_token, time.Now().UTC().Add(time.Hour))
	} else {
		_, err = db.Exec("UPDATE UserSessions SET SessionToken = ?, Expires_At = ? WHERE Username = ?", session_token, time.Now().UTC().Add(time.Hour), username)
	}

	return err
}

func UserSession_Check(username, session_token string) error {
	var exists int

	return db.QueryRow("SELECT 1 FROM UserSessions WHERE Username = ? AND SessionToken = ? AND ? < Expires_At", time.Now().UTC()).Scan(&exists)
}

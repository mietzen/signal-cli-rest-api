package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func GetDatabaseConnection() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		GetEnv("POSTGRES_HOST", ""),
		GetEnv("POSTGRES_PORT", ""),
		GetEnv("POSTGRES_USER", ""),
		GetEnv("POSTGRES_PASSWORD", ""),
		GetEnv("POSTGRES_DB", ""))

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	return db
}

func PushReceivedMessageToDatabase(message string) error {

	db := GetDatabaseConnection()

	insertStament := `insert into "received_messages"("data") values($1)`
	_, e := db.Exec(insertStament, message)
	defer db.Close()
	return e
}

func PushSendMessageToDatabase(message []byte) error {

	db := GetDatabaseConnection()

	insertStament := `insert into "send_messages"("data") values($1)`
	_, e := db.Exec(insertStament, message)
	defer db.Close()
	return e
}

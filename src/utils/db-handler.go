package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbConnection *sql.DB
)

func getDatabaseConnection() *sql.DB {
	if dbConnection == nil {
		dbType := GetEnv("DB_TYPE", "sqlite3")
		var dbSource string
		switch dbType {
		case "postgres":
			dbSource = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				GetEnv("DB_HOST", ""),
				GetEnv("DB_PORT", ""),
				GetEnv("DB_USER", ""),
				GetEnv("DB_PASSWORD", ""),
				GetEnv("DB_NAME", ""))
		default:
			// sqlite3
			dbSource = "/home/message-archive.db"
		}
		var err error
		dbConnection, err = sql.Open(dbType, dbSource)
		if err != nil {
			panic(err)
		}
	}
	return dbConnection
}

func getLastMsgs(numberOfMessages int, table string) (string, error) {
	if numberOfMessages == 0 {
		numberOfMessages = 5
	}
	var jsonStr string
	db := getDatabaseConnection()
	rows, err := db.Query(`SELECT "data" FROM "$1" ORDER BY id DESC LIMIT $2;`, table, numberOfMessages)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	jsonStr = "["
	i := 0
	for rows.Next() {
		var msg string
		err = rows.Scan(&msg)
		if err != nil {
			return "", err
		}
		jsonStr += msg
		if i != (numberOfMessages - 1) {
			jsonStr += ","
		}
		i++
	}
	jsonStr += "]"
	return jsonStr, err
}

func pushMsgsToDB(message string, table string) error {
	db := getDatabaseConnection()
	_, err := db.Exec(`insert into "$1"("data") values($2)`, table, message)
	return err
}

func PushReceivedMsgsToDB(message string) error {
	return pushMsgsToDB(message, "received_messages")
}

func PushSendMsgsToDB(message string) error {
	return pushMsgsToDB(message, "send_messages")
}

func GetLastReceivedMsgs(numberOfMessages int) (string, error) {
	return getLastMsgs(numberOfMessages, "received_messages")
}

func GetLastSendMsgs(numberOfMessages int) (string, error) {
	return getLastMsgs(numberOfMessages, "send_messages")
}

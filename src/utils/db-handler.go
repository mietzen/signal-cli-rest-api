package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

var (
	dbConnection *sql.DB
)

func getDatabaseConnection() *sql.DB {
	if dbConnection == nil {
		dbType := GetEnv("DB_TYPE", "sqlite")
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
			// sqlite
			dbSource = "/home/message-archive/message-archive.db"
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
	// Check message count
	queryStmt := fmt.Sprintf("SELECT COUNT(*) FROM %s;", table)
	rows, err := db.Query(queryStmt)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var n int
	err = rows.Scan(&n)
	if err != nil {
		return "", err
	}

	if numberOfMessages > n {
		numberOfMessages = n
	}
	if numberOfMessages != 0 {
		queryStmt = fmt.Sprintf("SELECT data FROM %s ORDER BY id DESC LIMIT %d;", table, numberOfMessages)
		rows, err = db.Query(queryStmt)
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
	}
	return jsonStr, err
}

func pushMsgsToDB(message string, table string) error {
	db := getDatabaseConnection()
	insertStmt := fmt.Sprintf("INSERT INTO %s (data) VALUES ('%s');", table, message)
	_, err := db.Exec(insertStmt)
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

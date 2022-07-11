package utils

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func GetIntEnv(key string, defaultVal int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		return intValue, nil
	}
	return defaultVal, nil
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsPhoneNumber(s string) bool {
	for index, c := range s {
		if index == 0 {
			if c != '+' {
				return false
			}
		} else {
			if c < '0' || c > '9' {
				return false
			}
		}
	}
	return true
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func GenerateAttachmentID() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_")
	id := make([]rune, 20)
	for i := range id {
		id[i] = letters[rand.Intn(len(letters))]
	}
	return string(id)
}

func GenerateAttachmentName(timestamp int64, mimetype string) string {
	t := time.Unix(timestamp, 0)
	return fmt.Sprintf("signal-%s.%s", t.Format("2006-01-02-150405"), mimetype)
}

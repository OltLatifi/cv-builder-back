package helpers

import (
	"os"
	"strconv"
)

func GetEnvInt(key string) int {
	value, _ := strconv.Atoi(os.Getenv(key))
	return value
}

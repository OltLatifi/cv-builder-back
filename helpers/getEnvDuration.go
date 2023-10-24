package helpers

import (
	"os"
	"time"
)

func GetEnvDuration(key string) (time.Duration, error) {
	return time.ParseDuration(os.Getenv(key))
}

package env

import (
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}

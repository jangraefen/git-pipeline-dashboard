package util

import (
	"os"
	"strconv"
	"strings"
)

func GetEnvOrDefault(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
}

func GetEnvOrDefaultUint(key string, val uint) uint {
	if val, ok := os.LookupEnv(key); ok {
		if parsedVal, err := strconv.ParseUint(val, 10, 32); err == nil {
			return uint(parsedVal)
		}
	}

	return val
}

func GetEnvList(key string) []string {
	if val := os.Getenv(key); val != "" {
		values := strings.Split(strings.TrimSpace(val), ",")
		for i, v := range values {
			values[i] = strings.TrimSpace(v)
		}

		return values
	}

	return []string{}
}

func ToEnvironmentName(s string) string {
	escaped := s

	escaped = strings.ToUpper(escaped)
	escaped = strings.ReplaceAll(escaped, "-", "_")
	escaped = strings.ReplaceAll(escaped, " ", "_")

	for strings.Contains(escaped, "__") {
		escaped = strings.ReplaceAll(escaped, "__", "_")
	}

	return escaped
}

package envparser

import (
	"fmt"
	"os"
	"strconv"
)

func GetConnectionStringPg() (string, error) {
	user, userExists := os.LookupEnv("POSTGRES_USER")
	password, passExists := os.LookupEnv("POSTGRES_PASSWORD")
	host, hostExists := os.LookupEnv("POSTGRES_HOST")
	portStr, portExists := os.LookupEnv("POSTGRES_PORT")
	dbName, dbExists := os.LookupEnv("POSTGRES_DB")
	sslMode, sslModeExists := os.LookupEnv("POSTGRES_SSLMODE")

	missingVars := []string{}
	if !userExists {
		missingVars = append(missingVars, "POSTGRES_USER")
	}
	if !passExists {
		missingVars = append(missingVars, "POSTGRES_PASSWORD")
	}
	if !hostExists {
		missingVars = append(missingVars, "POSTGRES_HOST")
	}
	if !portExists {
		missingVars = append(missingVars, "POSTGRES_PORT")
	}
	if !dbExists {
		missingVars = append(missingVars, "POSTGRES_DB")
	}
	if !sslModeExists {
		missingVars = append(missingVars, "POSTGRES_SSLMODE")
	}

	if len(missingVars) > 0 {
		return "", fmt.Errorf("the necessary environment variables are missing: %v", missingVars)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", fmt.Errorf("incorrect port number: %s", portStr)
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode,
	)

	return connStr, nil
}

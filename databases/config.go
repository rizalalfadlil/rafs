package databases

import (
	"fmt"
	"os"
)

// getDSN mengembalikan connection string PostgreSQL sesuai parameter credentials dari request
func getDSN(username, password, dbName string) string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "db"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, username, password, dbName, sslmode)
}

// getSuperadminDSN mengembalikan connection string PostgreSQL untuk superuser/database induk
func getSuperadminDSN() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "db"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_SUPERUSER")
	if user == "" {
		user = "superadmin"
	}
	pass := os.Getenv("DB_SUPERPASSWORD")
	if pass == "" {
		pass = "supersecret123"
	}
	dbname := os.Getenv("DB_SUPERDB")
	if dbname == "" {
		dbname = "postgres"
	}
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, dbname, sslmode)
}

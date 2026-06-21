package databases

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

var identifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// DBRequest adalah struktur data untuk menerima request JSON terkait database
type DBRequest struct {
	DBName   string `json:"db_name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func isValidIdentifier(name string) bool {
	return identifierRegex.MatchString(name)
}

// VerifyCredentials memverifikasi apakah username dan password cocok untuk mengakses database tertentu
func VerifyCredentials(dbName, username, password string) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database atau username tidak valid")
	}
	dsn := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("kredensial database tidak valid atau gagal terhubung: %w", err)
	}
	defer db.Close()
	return db.Ping()
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"status": "error", "message": message})
}

// CreateDatabaseHandlerCompat adalah endpoint legacy /api/create-db (POST only) untuk kompatibilitas ke belakang
func CreateDatabaseHandlerCompat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	var req DBRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.DBName == "" || req.Username == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Input tidak valid atau ada data yang kosong")
		return
	}

	err = CreateUserAndDatabase(req.DBName, req.Username, req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"status":  "sukses",
		"message": fmt.Sprintf("Database '%s' dan User '%s' berhasil dibuat!", req.DBName, req.Username),
	})
}

// DatabasesHandler mengelola CRUD database di endpoint /api/databases
func DatabasesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		dbs, err := ListDatabases()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"status":    "sukses",
			"databases": dbs,
		})

	case http.MethodPost:
		var req DBRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.DBName == "" || req.Username == "" || req.Password == "" {
			respondWithError(w, http.StatusBadRequest, "Input tidak valid atau ada data yang kosong")
			return
		}
		err = CreateUserAndDatabase(req.DBName, req.Username, req.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, map[string]string{
			"status":  "sukses",
			"message": fmt.Sprintf("Database '%s' dan User '%s' berhasil dibuat!", req.DBName, req.Username),
		})

	case http.MethodPut:
		username := r.Header.Get("X-Database-User")
		password := r.Header.Get("X-Database-Password")
		if username == "" || password == "" {
			respondWithError(w, http.StatusUnauthorized, "Autentikasi database diperlukan")
			return
		}

		var req struct {
			OldName string `json:"old_name"`
			NewName string `json:"new_name"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.OldName == "" || req.NewName == "" {
			respondWithError(w, http.StatusBadRequest, "Input tidak valid")
			return
		}

		err = VerifyCredentials(req.OldName, username, password)
		if err != nil {
			respondWithError(w, http.StatusForbidden, "Autentikasi database gagal: "+err.Error())
			return
		}

		err = RenameDatabase(req.OldName, req.NewName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "sukses",
			"message": fmt.Sprintf("Database '%s' berhasil diubah nama menjadi '%s'!", req.OldName, req.NewName),
		})

	case http.MethodDelete:
		username := r.Header.Get("X-Database-User")
		password := r.Header.Get("X-Database-Password")
		if username == "" || password == "" {
			respondWithError(w, http.StatusUnauthorized, "Autentikasi database diperlukan")
			return
		}

		dbName := r.URL.Query().Get("db_name")
		if dbName == "" {
			respondWithError(w, http.StatusBadRequest, "Parameter db_name wajib disertakan")
			return
		}

		err := VerifyCredentials(dbName, username, password)
		if err != nil {
			respondWithError(w, http.StatusForbidden, "Autentikasi database gagal: "+err.Error())
			return
		}

		err = DeleteDatabase(dbName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "sukses",
			"message": fmt.Sprintf("Database '%s' berhasil dihapus!", dbName),
		})

	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
	}
}

// CreateUserAndDatabase mengelola pembuatan user baru dan database baru di PostgreSQL
func CreateUserAndDatabase(dbName, username, password string) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database atau username tidak valid")
	}
	dsn := "host=db port=5432 user=superadmin password=supersecret123 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database induk: %w", err)
	}
	defer db.Close()

	createUserSQL := fmt.Sprintf("CREATE USER %s WITH PASSWORD '%s';", username, password)
	_, err = db.Exec(createUserSQL)
	if err != nil {
		return fmt.Errorf("gagal membuat user: %w", err)
	}

	createDBSQL := fmt.Sprintf("CREATE DATABASE %s OWNER %s;", dbName, username)
	_, err = db.Exec(createDBSQL)
	if err != nil {
		db.Exec(fmt.Sprintf("DROP USER %s;", username))
		return fmt.Errorf("gagal membuat database: %w", err)
	}

	return nil
}

// ListDatabases mengembalikan daftar database buatan pengguna
func ListDatabases() ([]string, error) {
	dsn := "host=db port=5432 user=superadmin password=supersecret123 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database induk: %w", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false AND datname NOT IN ('postgres', 'db');")
	if err != nil {
		return nil, fmt.Errorf("gagal query database: %w", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		databases = append(databases, name)
	}
	return databases, nil
}

// RenameDatabase mengubah nama database
func RenameDatabase(oldName, newName string) error {
	if !isValidIdentifier(oldName) || !isValidIdentifier(newName) {
		return fmt.Errorf("nama database tidak valid")
	}

	dsn := "host=db port=5432 user=superadmin password=supersecret123 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("ALTER DATABASE %s RENAME TO %s;", oldName, newName)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("gagal mengubah nama database: %w", err)
	}
	return nil
}

// DeleteDatabase menghapus database beserta koneksinya
func DeleteDatabase(dbName string) error {
	if !isValidIdentifier(dbName) {
		return fmt.Errorf("nama database tidak valid")
	}

	dsn := "host=db port=5432 user=superadmin password=supersecret123 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	disconnectQuery := fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%s'
		  AND pid <> pg_backend_pid();`, dbName)
	_, _ = db.Exec(disconnectQuery)

	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbName)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("gagal menghapus database: %w", err)
	}
	return nil
}

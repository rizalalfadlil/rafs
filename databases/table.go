package databases

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

var typeRegex = regexp.MustCompile(`^[a-zA-Z]+(?:\([0-9]+\))?(?:\s+[a-zA-Z0-9_\s]+)*$`)

// ColumnDef menyimpan definisi kolom untuk pembuatan tabel
type ColumnDef struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func isValidType(t string) bool {
	return typeRegex.MatchString(t)
}

func joinStrings(slice []string, sep string) string {
	if len(slice) == 0 {
		return ""
	}
	res := slice[0]
	for _, s := range slice[1:] {
		res += sep + s
	}
	return res
}

// TablesHandler mengelola CRUD tabel di dalam database tertentu
func TablesHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Database-User")
	password := r.Header.Get("X-Database-Password")
	if username == "" || password == "" {
		respondWithError(w, http.StatusUnauthorized, "Autentikasi database diperlukan")
		return
	}

	switch r.Method {
	case http.MethodGet:
		dbName := r.URL.Query().Get("db_name")
		if dbName == "" {
			respondWithError(w, http.StatusBadRequest, "Parameter db_name wajib disertakan")
			return
		}
		tables, err := ListTables(dbName, username, password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"status": "sukses",
			"tables": tables,
		})

	case http.MethodPost:
		var req struct {
			DBName    string      `json:"db_name"`
			TableName string      `json:"table_name"`
			Columns   []ColumnDef `json:"columns"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.DBName == "" || req.TableName == "" || len(req.Columns) == 0 {
			respondWithError(w, http.StatusBadRequest, "Input tidak valid")
			return
		}
		err = CreateTable(req.DBName, username, password, req.TableName, req.Columns)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, map[string]string{
			"status":  "sukses",
			"message": fmt.Sprintf("Tabel '%s' berhasil dibuat di database '%s'!", req.TableName, req.DBName),
		})

	case http.MethodPut:
		var req struct {
			DBName  string `json:"db_name"`
			OldName string `json:"old_name"`
			NewName string `json:"new_name"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.DBName == "" || req.OldName == "" || req.NewName == "" {
			respondWithError(w, http.StatusBadRequest, "Input tidak valid")
			return
		}
		err = RenameTable(req.DBName, username, password, req.OldName, req.NewName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "sukses",
			"message": fmt.Sprintf("Tabel '%s' berhasil diubah nama menjadi '%s' di database '%s'!", req.OldName, req.NewName, req.DBName),
		})

	case http.MethodDelete:
		dbName := r.URL.Query().Get("db_name")
		tableName := r.URL.Query().Get("table_name")
		if dbName == "" || tableName == "" {
			respondWithError(w, http.StatusBadRequest, "Parameter db_name dan table_name wajib disertakan")
			return
		}
		err := DeleteTable(dbName, username, password, tableName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "sukses",
			"message": fmt.Sprintf("Tabel '%s' berhasil dihapus dari database '%s'!", tableName, dbName),
		})

	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
	}
}

// ListTables mengembalikan daftar tabel di database tertentu
func ListTables(dbName, username, password string) ([]string, error) {
	if !isValidIdentifier(dbName) || !isValidIdentifier(username) {
		return nil, fmt.Errorf("nama database atau username tidak valid")
	}

	dsn := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database %s: %w", dbName, err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';")
	if err != nil {
		return nil, fmt.Errorf("gagal query tabel: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, nil
}

// CreateTable membuat tabel baru di database tertentu
func CreateTable(dbName, username, password, tableName string, columns []ColumnDef) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database, username, atau tabel tidak valid")
	}
	if len(columns) == 0 {
		return fmt.Errorf("kolom tidak boleh kosong")
	}

	var colQueries []string
	for _, col := range columns {
		if !isValidIdentifier(col.Name) {
			return fmt.Errorf("nama kolom '%s' tidak valid", col.Name)
		}
		if !isValidType(col.Type) {
			return fmt.Errorf("tipe data '%s' tidak valid", col.Type)
		}
		colQueries = append(colQueries, fmt.Sprintf("%s %s", col.Name, col.Type))
	}

	dsn := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("CREATE TABLE %s (%s);", tableName, joinStrings(colQueries, ", "))
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("gagal membuat tabel: %w", err)
	}
	return nil
}

// RenameTable mengubah nama tabel di database tertentu
func RenameTable(dbName, username, password, oldName, newName string) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(oldName) || !isValidIdentifier(newName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database, username, atau tabel tidak valid")
	}

	dsn := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", oldName, newName)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("gagal mengubah nama tabel: %w", err)
	}
	return nil
}

// DeleteTable menghapus tabel di database tertentu
func DeleteTable(dbName, username, password, tableName string) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database, username, atau tabel tidak valid")
	}

	dsn := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("gagal menghapus tabel: %w", err)
	}
	return nil
}

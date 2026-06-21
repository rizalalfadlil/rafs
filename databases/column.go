package databases

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// ListColumns mengembalikan daftar kolom beserta tipe datanya untuk tabel tertentu
func ListColumns(dbName, username, password, tableName string) ([]map[string]string, error) {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(username) {
		return nil, fmt.Errorf("nama database, username, atau tabel tidak valid")
	}

	dsn := getDSN(username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database %s: %w", dbName, err)
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT column_name, data_type, is_nullable, column_default 
		FROM information_schema.columns 
		WHERE table_schema = 'public' AND table_name = $1 
		ORDER BY ordinal_position;`, tableName)
	if err != nil {
		return nil, fmt.Errorf("gagal query kolom: %w", err)
	}
	defer rows.Close()

	var columns []map[string]string
	for rows.Next() {
		var name, dataType, isNullable string
		var columnDefault sql.NullString
		if err := rows.Scan(&name, &dataType, &isNullable, &columnDefault); err != nil {
			return nil, err
		}
		
		defaultVal := ""
		if columnDefault.Valid {
			defaultVal = columnDefault.String
		}

		columns = append(columns, map[string]string{
			"name":     name,
			"type":     dataType,
			"nullable": isNullable,
			"default":  defaultVal,
		})
	}
	return columns, nil
}

// AddColumn menambahkan kolom baru ke dalam tabel
func AddColumn(dbName, username, password, tableName, colName, colType string) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(colName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database, username, tabel, atau kolom tidak valid")
	}
	if !isValidType(colType) {
		return fmt.Errorf("tipe data tidak valid")
	}

	dsn := getDSN(username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", tableName, colName, colType)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("gagal menambahkan kolom: %w", err)
	}
	return nil
}

// DeleteColumn menghapus kolom dari tabel
func DeleteColumn(dbName, username, password, tableName, colName string) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(colName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database, username, tabel, atau kolom tidak valid")
	}

	dsn := getDSN(username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", tableName, colName)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("gagal menghapus kolom: %w", err)
	}
	return nil
}

// RenameColumn mengubah nama kolom di tabel
func RenameColumn(dbName, username, password, tableName, oldName, newName string) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(oldName) || !isValidIdentifier(newName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database, username, tabel, atau kolom tidak valid")
	}

	dsn := getDSN(username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("ALTER TABLE %s RENAME COLUMN %s TO %s;", tableName, oldName, newName)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("gagal mengubah nama kolom: %w", err)
	}
	return nil
}

// AlterColumnType mengubah tipe data kolom di tabel
func AlterColumnType(dbName, username, password, tableName, colName, newType string) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(colName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database, username, tabel, atau kolom tidak valid")
	}
	if !isValidType(newType) {
		return fmt.Errorf("tipe data tidak valid")
	}

	dsn := getDSN(username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN %s TYPE %s;", tableName, colName, newType)
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("gagal mengubah tipe data kolom: %w", err)
	}
	return nil
}

// ColumnsHandler mengelola request CRUD kolom
func ColumnsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Database-User")
	password := r.Header.Get("X-Database-Password")
	if username == "" || password == "" {
		respondWithError(w, http.StatusUnauthorized, "Autentikasi database diperlukan")
		return
	}

	switch r.Method {
	case http.MethodGet:
		dbName := r.URL.Query().Get("db_name")
		tableName := r.URL.Query().Get("table_name")
		if dbName == "" || tableName == "" {
			respondWithError(w, http.StatusBadRequest, "Parameter db_name dan table_name wajib disertakan")
			return
		}
		cols, err := ListColumns(dbName, username, password, tableName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"status":  "sukses",
			"columns": cols,
		})

	case http.MethodPost:
		var req struct {
			DBName     string `json:"db_name"`
			TableName  string `json:"table_name"`
			ColumnName string `json:"column_name"`
			ColumnType string `json:"column_type"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.DBName == "" || req.TableName == "" || req.ColumnName == "" || req.ColumnType == "" {
			respondWithError(w, http.StatusBadRequest, "Input tidak valid atau ada data yang kosong")
			return
		}
		err = AddColumn(req.DBName, username, password, req.TableName, req.ColumnName, req.ColumnType)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, map[string]string{
			"status":  "sukses",
			"message": fmt.Sprintf("Kolom '%s' berhasil ditambahkan ke tabel '%s'!", req.ColumnName, req.TableName),
		})

	case http.MethodPut:
		var req struct {
			DBName     string `json:"db_name"`
			TableName  string `json:"table_name"`
			OldName    string `json:"old_name"`
			NewName    string `json:"new_name"`
			ColumnType string `json:"column_type"` // optional
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.DBName == "" || req.TableName == "" || req.OldName == "" || req.NewName == "" {
			respondWithError(w, http.StatusBadRequest, "Input tidak valid")
			return
		}
		if req.OldName != req.NewName {
			err = RenameColumn(req.DBName, username, password, req.TableName, req.OldName, req.NewName)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		if req.ColumnType != "" {
			err = AlterColumnType(req.DBName, username, password, req.TableName, req.NewName, req.ColumnType)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "sukses",
			"message": fmt.Sprintf("Kolom '%s' berhasil diubah!", req.OldName),
		})

	case http.MethodDelete:
		dbName := r.URL.Query().Get("db_name")
		tableName := r.URL.Query().Get("table_name")
		columnName := r.URL.Query().Get("column_name")
		if dbName == "" || tableName == "" || columnName == "" {
			respondWithError(w, http.StatusBadRequest, "Parameter db_name, table_name, dan column_name wajib disertakan")
			return
		}
		err := DeleteColumn(dbName, username, password, tableName, columnName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "sukses",
			"message": fmt.Sprintf("Kolom '%s' berhasil dihapus dari tabel '%s'!", columnName, tableName),
		})

	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
	}
}

package databases

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// ListRows mengembalikan semua baris data dari tabel tertentu
func ListRows(dbName, username, password, tableName string) ([]map[string]interface{}, error) {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(username) {
		return nil, fmt.Errorf("nama database, username, atau tabel tidak valid")
	}

	dsn := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database %s: %w", dbName, err)
	}
	defer db.Close()

	queryStr := fmt.Sprintf("SELECT * FROM %s;", tableName)
	rows, err := db.Query(queryStr)
	if err != nil {
		return nil, fmt.Errorf("gagal query baris: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil nama kolom: %w", err)
	}

	var allRows []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("gagal scan baris: %w", err)
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columns[i]
			b, ok := val.([]byte)
			if ok {
				m[colName] = string(b)
			} else {
				m[colName] = val
			}
		}
		allRows = append(allRows, m)
	}

	// Jika nil (tabel kosong), kembalikan array kosong agar json-encode menghasilkan [] bukan null
	if allRows == nil {
		allRows = []map[string]interface{}{}
	}

	return allRows, nil
}

// InsertRow menyisipkan baris baru ke dalam tabel
func InsertRow(dbName, username, password, tableName string, rowData map[string]interface{}) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database, username, atau tabel tidak valid")
	}
	if len(rowData) == 0 {
		return fmt.Errorf("data baris tidak boleh kosong")
	}

	var cols []string
	var placeholders []string
	var vals []interface{}
	i := 1
	for colName, colVal := range rowData {
		if !isValidIdentifier(colName) {
			return fmt.Errorf("nama kolom '%s' tidak valid", colName)
		}
		cols = append(cols, colName)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		vals = append(vals, colVal)
		i++
	}

	dsn := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", tableName, joinStrings(cols, ", "), joinStrings(placeholders, ", "))
	_, err = db.Exec(query, vals...)
	if err != nil {
		return fmt.Errorf("gagal menyisipkan baris: %w", err)
	}
	return nil
}

// UpdateRow memperbarui baris di tabel berdasarkan filter kriteria
func UpdateRow(dbName, username, password, tableName string, whereData, rowData map[string]interface{}) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database, username, atau tabel tidak valid")
	}
	if len(rowData) == 0 {
		return fmt.Errorf("data set tidak boleh kosong")
	}
	if len(whereData) == 0 {
		return fmt.Errorf("kriteria update (where) tidak boleh kosong")
	}

	var setExprs []string
	var whereExprs []string
	var vals []interface{}
	i := 1
	for colName, colVal := range rowData {
		if !isValidIdentifier(colName) {
			return fmt.Errorf("nama kolom set '%s' tidak valid", colName)
		}
		setExprs = append(setExprs, fmt.Sprintf("%s = $%d", colName, i))
		vals = append(vals, colVal)
		i++
	}
	for colName, colVal := range whereData {
		if !isValidIdentifier(colName) {
			return fmt.Errorf("nama kolom kriteria '%s' tidak valid", colName)
		}
		whereExprs = append(whereExprs, fmt.Sprintf("%s = $%d", colName, i))
		vals = append(vals, colVal)
		i++
	}

	dsn := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s;", tableName, joinStrings(setExprs, ", "), joinStrings(whereExprs, " AND "))
	_, err = db.Exec(query, vals...)
	if err != nil {
		return fmt.Errorf("gagal memperbarui baris: %w", err)
	}
	return nil
}

// DeleteRow menghapus baris dari tabel berdasarkan kriteria filter
func DeleteRow(dbName, username, password, tableName string, whereData map[string]interface{}) error {
	if !isValidIdentifier(dbName) || !isValidIdentifier(tableName) || !isValidIdentifier(username) {
		return fmt.Errorf("nama database, username, atau tabel tidak valid")
	}
	if len(whereData) == 0 {
		return fmt.Errorf("kriteria hapus (where) tidak boleh kosong")
	}

	var whereExprs []string
	var vals []interface{}
	i := 1
	for colName, colVal := range whereData {
		if !isValidIdentifier(colName) {
			return fmt.Errorf("nama kolom kriteria '%s' tidak valid", colName)
		}
		whereExprs = append(whereExprs, fmt.Sprintf("%s = $%d", colName, i))
		vals = append(vals, colVal)
		i++
	}

	dsn := fmt.Sprintf("host=db port=5432 user=%s password=%s dbname=%s sslmode=disable", username, password, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke database: %w", err)
	}
	defer db.Close()

	query := fmt.Sprintf("DELETE FROM %s WHERE %s;", tableName, joinStrings(whereExprs, " AND "))
	_, err = db.Exec(query, vals...)
	if err != nil {
		return fmt.Errorf("gagal menghapus baris: %w", err)
	}
	return nil
}

// RowsHandler mengelola request CRUD baris
func RowsHandler(w http.ResponseWriter, r *http.Request) {
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
		rows, err := ListRows(dbName, username, password, tableName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"status": "sukses",
			"rows":   rows,
		})

	case http.MethodPost:
		var req struct {
			DBName    string                 `json:"db_name"`
			TableName string                 `json:"table_name"`
			Row       map[string]interface{} `json:"row"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.DBName == "" || req.TableName == "" || len(req.Row) == 0 {
			respondWithError(w, http.StatusBadRequest, "Input tidak valid atau data kosong")
			return
		}
		err = InsertRow(req.DBName, username, password, req.TableName, req.Row)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusCreated, map[string]string{
			"status":  "sukses",
			"message": "Baris baru berhasil ditambahkan!",
		})

	case http.MethodPut:
		var req struct {
			DBName    string                 `json:"db_name"`
			TableName string                 `json:"table_name"`
			Where     map[string]interface{} `json:"where"`
			Row       map[string]interface{} `json:"row"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.DBName == "" || req.TableName == "" || len(req.Where) == 0 || len(req.Row) == 0 {
			respondWithError(w, http.StatusBadRequest, "Input tidak valid atau data kosong")
			return
		}
		err = UpdateRow(req.DBName, username, password, req.TableName, req.Where, req.Row)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "sukses",
			"message": "Baris berhasil diperbarui!",
		})

	case http.MethodDelete:
		var req struct {
			DBName    string                 `json:"db_name"`
			TableName string                 `json:"table_name"`
			Where     map[string]interface{} `json:"where"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil || req.DBName == "" || req.TableName == "" || len(req.Where) == 0 {
			respondWithError(w, http.StatusBadRequest, "Input tidak valid atau kriteria kosong")
			return
		}
		err = DeleteRow(req.DBName, username, password, req.TableName, req.Where)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "sukses",
			"message": "Baris berhasil dihapus!",
		})

	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
	}
}

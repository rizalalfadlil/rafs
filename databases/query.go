package databases

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

// QueryRequest mewakili data input untuk menjalankan perintah SQL
type QueryRequest struct {
	DBName string `json:"db_name"`
	Query  string `json:"query"`
}

// QueryHandler menerima input query SQL dari admin panel dan mengeksekusinya
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	var req QueryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.DBName == "" || req.Query == "" {
		respondWithError(w, http.StatusBadRequest, "Input tidak valid atau ada data yang kosong")
		return
	}

	username := r.Header.Get("X-Database-User")
	password := r.Header.Get("X-Database-Password")
	if username == "" || password == "" {
		respondWithError(w, http.StatusUnauthorized, "Autentikasi database diperlukan (Header X-Database-User & X-Database-Password kosong)")
		return
	}

	// Hubungkan ke database dengan kredensial yang diberikan
	dsn := getDSN(username, password, req.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal membuka koneksi database: "+err.Error())
		return
	}
	defer db.Close()

	// Ping untuk memastikan kredensial benar
	err = db.Ping()
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Autentikasi database gagal: "+err.Error())
		return
	}

	trimmedQuery := strings.ToLower(strings.TrimSpace(req.Query))
	isSelect := strings.HasPrefix(trimmedQuery, "select") || strings.HasPrefix(trimmedQuery, "show") || strings.HasPrefix(trimmedQuery, "explain")

	if isSelect {
		// Jalankan SELECT query
		rows, err := db.Query(req.Query)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Gagal mengambil metadata kolom: "+err.Error())
			return
		}

		var allRows []map[string]interface{}
		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				respondWithError(w, http.StatusInternalServerError, "Gagal membaca baris data: "+err.Error())
				return
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

		if allRows == nil {
			allRows = []map[string]interface{}{}
		}

		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"status":  "sukses",
			"type":    "select",
			"columns": cols,
			"rows":    allRows,
		})
	} else {
		// Jalankan perintah DDL/DML (INSERT, UPDATE, DELETE, CREATE, DROP, etc.)
		res, err := db.Exec(req.Query)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		rowsAffected, _ := res.RowsAffected()
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"status":        "sukses",
			"type":          "exec",
			"rows_affected": rowsAffected,
			"message":       "Perintah berhasil dieksekusi.",
		})
	}
}

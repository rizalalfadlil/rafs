package main

import (
	"fmt"
	"log"
	"net/http"

	"rafs/databases"
	"rafs/sites"
)

func main() {
	// Jalankan file server untuk web statis (dari langkah sebelumnya)
	fileServer := http.FileServer(http.Dir("./www"))
	http.Handle("/sites/", http.StripPrefix("/sites/", fileServer))

	// Jalankan file server untuk admin panel (Vue app)
	adminServer := http.FileServer(http.Dir("./admin/dist"))
	http.Handle("/admin/", http.StripPrefix("/admin/", adminServer))

	// Endpoint API baru untuk mengelola database dari package databases
	http.HandleFunc("/api/create-db", databases.CreateDatabaseHandlerCompat)
	http.HandleFunc("/api/databases", databases.DatabasesHandler)
	http.HandleFunc("/api/tables", databases.TablesHandler)
	http.HandleFunc("/api/columns", databases.ColumnsHandler)
	http.HandleFunc("/api/rows", databases.RowsHandler)
	http.HandleFunc("/api/server-info", databases.ServerInfoHandler)

	// Endpoint API baru untuk mengelola website statis (sites)
	http.HandleFunc("/api/sites", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sites.ListSitesHandler(w, r)
		case http.MethodDelete:
			sites.DeleteSiteHandler(w, r)
		default:
			w.Header().Set("Allow", "GET, DELETE")
			http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/sites/clone", sites.CloneSiteHandler)
	http.HandleFunc("/api/sites/upload", sites.UploadSiteHandler)

	// Rute utama - dikembalikan ke default awal
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Server Utama Aktif. Gunakan API /api/databases dan /api/tables untuk mengelola database.")
	})

	port := ":8080"
	fmt.Printf("Server berjalan di port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

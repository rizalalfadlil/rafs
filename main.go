package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"rafs/databases"
	"rafs/sites"
	"rafs/storage"
)

func main() {
	// Buat folder-folder yang dibutuhkan jika belum ada
	_ = os.MkdirAll("./www", 0755)
	_ = os.MkdirAll("./storage_data", 0755)
	_ = os.MkdirAll("./public", 0755)

	// Jalankan file server untuk web statis (dari langkah sebelumnya)
	fileServer := http.FileServer(http.Dir("./www"))
	http.Handle("/sites/", http.StripPrefix("/sites/", fileServer))

	// Jalankan file server untuk berkas publik
	publicServer := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public/", publicServer))

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
	http.HandleFunc("/api/query", databases.QueryHandler)

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

	// Endpoint API baru untuk mengelola cloud storage
	http.HandleFunc("/api/storage", storage.ListStorageHandler)
	http.HandleFunc("/api/storage/folder", storage.CreateFolderHandler)
	http.HandleFunc("/api/storage/upload", storage.UploadFilesHandler)
	http.HandleFunc("/api/storage/delete", storage.DeleteItemsHandler)
	http.HandleFunc("/api/storage/download", storage.DownloadFileHandler)
	http.HandleFunc("/api/storage/public", storage.SetPublicHandler)

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

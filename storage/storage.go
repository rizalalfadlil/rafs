package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var safeNameRegex = regexp.MustCompile(`^[a-zA-Z0-9 _.-]+$`)

const storageRoot = "./storage_data"
const publicRoot = "./public"

func isValidItemName(name string) bool {
	return safeNameRegex.MatchString(name)
}

// ItemInfo represents a file or folder in storage
type ItemInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"` // "file" or "folder"
	Size    int64  `json:"size,omitempty"`
	ModTime string `json:"mod_time,omitempty"`
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"status": "error", "message": message})
}

// getSafePath validates path to prevent directory traversal
func getSafePath(baseDir string, subPath string) (string, error) {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}

	joined := filepath.Join(absBase, filepath.Clean(subPath))
	absJoined, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}

	if absJoined != absBase && !strings.HasPrefix(absJoined, absBase+string(filepath.Separator)) {
		return "", fmt.Errorf("akses di luar batas direktori aman ditolak")
	}

	return absJoined, nil
}

// ListStorageHandler lists files and folders in ./storage/<path>
func ListStorageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	subPath := r.URL.Query().Get("path")
	storageDir, err := getSafePath(storageRoot, subPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	entries, err := os.ReadDir(storageDir)
	if err != nil {
		if os.IsNotExist(err) {
			respondWithError(w, http.StatusNotFound, "Direktori tidak ditemukan")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Gagal membaca direktori: "+err.Error())
		return
	}

	var contents []ItemInfo
	for _, entry := range entries {
		info, err := entry.Info()
		var size int64
		var modTime string
		if err == nil {
			size = info.Size()
			modTime = info.ModTime().Format("2006-01-02 15:04:05")
		}

		itemType := "file"
		if entry.IsDir() {
			itemType = "folder"
		}

		contents = append(contents, ItemInfo{
			Name:    entry.Name(),
			Type:    itemType,
			Size:    size,
			ModTime: modTime,
		})
	}

	if contents == nil {
		contents = []ItemInfo{}
	}

	// Hitung total penggunaan storage
	var totalUsed int64
	filepath.Walk(storageRoot, func(_ string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			totalUsed += info.Size()
		}
		return nil
	})

	// Quota fixed pada 1GB
	const totalQuota int64 = 1024 * 1024 * 1024

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":       "sukses",
		"current_path": filepath.ToSlash(subPath),
		"contents":     contents,
		"space": map[string]int64{
			"used":  totalUsed,
			"total": totalQuota,
		},
	})
}

// CreateFolderHandler creates a folder under ./storage/<path>/<name>
func CreateFolderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	var req struct {
		Path string `json:"path"`
		Name string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Input tidak valid")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if !isValidItemName(req.Name) {
		respondWithError(w, http.StatusBadRequest, "Nama folder mengandung karakter yang dilarang")
		return
	}

	parentDir, err := getSafePath(storageRoot, req.Path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	targetDir := filepath.Join(parentDir, req.Name)
	if _, err := os.Stat(targetDir); err == nil {
		respondWithError(w, http.StatusBadRequest, "Folder sudah ada")
		return
	}

	err = os.Mkdir(targetDir, 0755)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal membuat folder: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{
		"status":  "sukses",
		"message": fmt.Sprintf("Folder '%s' berhasil dibuat!", req.Name),
	})
}

// UploadFilesHandler receives file uploads to ./storage/<path>
func UploadFilesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	err := r.ParseMultipartForm(100 << 20) // Maks 100MB
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Gagal memproses upload")
		return
	}

	targetPath := r.FormValue("path")
	destDir, err := getSafePath(storageRoot, targetPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		respondWithError(w, http.StatusBadRequest, "Tidak ada file yang diunggah")
		return
	}

	uploadedCount := 0
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Bersihkan nama file
		safeName := filepath.Base(fileHeader.Filename)
		if !isValidItemName(safeName) {
			safeName = "file_upload"
		}

		targetFilePath := filepath.Join(destDir, safeName)

		out, err := os.Create(targetFilePath)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Gagal menyimpan file: "+err.Error())
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Gagal menulis file: "+err.Error())
			return
		}
		uploadedCount++
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "sukses",
		"message": fmt.Sprintf("%d berkas berhasil diunggah!", uploadedCount),
	})
}

// DeleteItemsHandler deletes files or folders recursively from ./storage/
func DeleteItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	var req struct {
		Paths []string `json:"paths"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.Paths) == 0 {
		respondWithError(w, http.StatusBadRequest, "Input tidak valid")
		return
	}

	absRoot, err := filepath.Abs(storageRoot)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal menginisialisasi root storage")
		return
	}

	deletedCount := 0
	for _, itemPath := range req.Paths {
		safePath, err := getSafePath(storageRoot, itemPath)
		if err != nil {
			continue
		}

		// Keamanan: jangan hapus root direktori
		if safePath == absRoot {
			continue
		}

		err = os.RemoveAll(safePath)
		if err == nil {
			deletedCount++
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "sukses",
		"message": fmt.Sprintf("%d file/folder berhasil dihapus!", deletedCount),
	})
}

// DownloadFileHandler serves file content for download/preview
func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	itemPath := r.URL.Query().Get("path")
	safePath, err := getSafePath(storageRoot, itemPath)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	info, err := os.Stat(safePath)
	if err != nil {
		if os.IsNotExist(err) {
			respondWithError(w, http.StatusNotFound, "File tidak ditemukan")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Gagal mengakses file: "+err.Error())
		return
	}

	if info.IsDir() {
		respondWithError(w, http.StatusBadRequest, "Tidak dapat mendownload direktori secara langsung")
		return
	}

	// Serve file secara native
	http.ServeFile(w, r, safePath)
}

// SetPublicHandler copies selected files to public folder and returns public URLs
func SetPublicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	var req struct {
		Paths []string `json:"paths"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.Paths) == 0 {
		respondWithError(w, http.StatusBadRequest, "Input tidak valid")
		return
	}

	absStorage, err := filepath.Abs(storageRoot)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal menginisialisasi root storage")
		return
	}

	var publicLinks []string
	for _, itemPath := range req.Paths {
		safePath, err := getSafePath(storageRoot, itemPath)
		if err != nil {
			continue
		}

		rel, err := filepath.Rel(absStorage, safePath)
		if err != nil {
			continue
		}

		destPublicPath := filepath.Join(publicRoot, rel)

		// Buat subfolder jika diperlukan
		err = os.MkdirAll(filepath.Dir(destPublicPath), 0755)
		if err != nil {
			continue
		}

		err = copyRecursive(safePath, destPublicPath)
		if err != nil {
			continue
		}

		// Menghasilkan link publik (sesuai base URL server)
		publicLinks = append(publicLinks, "/public/"+filepath.ToSlash(rel))
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "sukses",
		"message": fmt.Sprintf("%d file/folder berhasil dipublikasikan!", len(publicLinks)),
		"links":   publicLinks,
	})
}

// Helper: Salin berkas atau direktori secara rekursif
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}

func copyRecursive(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		err = os.MkdirAll(dst, srcInfo.Mode())
		if err != nil {
			return err
		}
		entries, err := os.ReadDir(src)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			err = copyRecursive(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name()))
			if err != nil {
				return err
			}
		}
	} else {
		err = copyFile(src, dst)
		if err != nil {
			return err
		}
	}
	return nil
}

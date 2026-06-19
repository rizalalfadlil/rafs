package sites

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var identifierRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func isValidSiteName(name string) bool {
	return identifierRegex.MatchString(name)
}

// SiteInfo menyimpan nama website dan status keaktifan
type SiteInfo struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"status": "error", "message": message})
}

// ListSitesHandler mengembalikan daftar semua website statis di folder ./www
func ListSitesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	wwwDir := "./www"
	entries, err := os.ReadDir(wwwDir)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal membaca direktori www: "+err.Error())
		return
	}

	var sites []SiteInfo
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			// Cek apakah ada index.html di dalam folder tersebut
			indexPath := filepath.Join(wwwDir, name, "index.html")
			active := false
			if _, err := os.Stat(indexPath); err == nil {
				active = true
			}
			sites = append(sites, SiteInfo{
				Name:   name,
				Active: active,
			})
		}
	}

	if sites == nil {
		sites = []SiteInfo{}
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": "sukses",
		"sites":  sites,
	})
}

// CloneSiteHandler melakukan clone repository github publik ke dalam folder ./www/<site_name>
func CloneSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	var req struct {
		RepoURL  string `json:"repo_url"`
		SiteName string `json:"site_name"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RepoURL == "" || req.SiteName == "" {
		respondWithError(w, http.StatusBadRequest, "Input tidak valid")
		return
	}

	req.SiteName = strings.TrimSpace(req.SiteName)
	if !isValidSiteName(req.SiteName) {
		respondWithError(w, http.StatusBadRequest, "Nama site tidak valid (hanya alfanumerik, dash, dan underscore)")
		return
	}

	targetDir := filepath.Join("./www", req.SiteName)
	if _, err := os.Stat(targetDir); err == nil {
		respondWithError(w, http.StatusBadRequest, "Folder website dengan nama tersebut sudah ada")
		return
	}

	// Jalankan perintah clone
	cmd := exec.Command("git", "clone", "--depth", "1", req.RepoURL, targetDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(targetDir)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Gagal clone repo: %s (Output: %s)", err.Error(), string(output)))
		return
	}

	// Hapus folder .git agar tersimpan bersih
	os.RemoveAll(filepath.Join(targetDir, ".git"))

	// Cek keaktifan (apakah ada index.html)
	indexPath := filepath.Join(targetDir, "index.html")
	active := false
	if _, err := os.Stat(indexPath); err == nil {
		active = true
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "sukses",
		"message": fmt.Sprintf("Repository berhasil diclone ke website '%s'!", req.SiteName),
		"active":  active,
	})
}

// UploadSiteHandler menerima upload file zip dan mengekstraknya di ./www/<site_name>
func UploadSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	err := r.ParseMultipartForm(50 << 20) // maks 50MB
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Gagal memproses form upload")
		return
	}

	siteName := r.FormValue("site_name")
	siteName = strings.TrimSpace(siteName)
	if !isValidSiteName(siteName) {
		respondWithError(w, http.StatusBadRequest, "Nama site tidak valid")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "File zip wajib diunggah")
		return
	}
	defer file.Close()

	targetDir := filepath.Join("./www", siteName)
	if _, err := os.Stat(targetDir); err == nil {
		respondWithError(w, http.StatusBadRequest, "Folder website dengan nama tersebut sudah ada")
		return
	}

	// Buat file zip sementara
	tempZipPath := filepath.Join(os.TempDir(), fmt.Sprintf("upload-%d.zip", time.Now().UnixNano()))
	tempFile, err := os.Create(tempZipPath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal membuat file sementara")
		return
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempZipPath)
	}()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal menyimpan file zip")
		return
	}
	tempFile.Close()

	// Ekstrak zip ke folder target
	err = unzip(tempZipPath, targetDir)
	if err != nil {
		os.RemoveAll(targetDir)
		respondWithError(w, http.StatusInternalServerError, "Gagal mengekstrak zip: "+err.Error())
		return
	}

	// Cek keaktifan
	indexPath := filepath.Join(targetDir, "index.html")
	active := false
	if _, err := os.Stat(indexPath); err == nil {
		active = true
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"status":  "sukses",
		"message": fmt.Sprintf("Website '%s' berhasil diunggah!", siteName),
		"active":  active,
	})
}

// DeleteSiteHandler menghapus folder website dari ./www
func DeleteSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondWithError(w, http.StatusMethodNotAllowed, "Metode tidak diizinkan")
		return
	}

	siteName := r.URL.Query().Get("site_name")
	if siteName == "" {
		respondWithError(w, http.StatusBadRequest, "Parameter site_name wajib disertakan")
		return
	}

	if !isValidSiteName(siteName) {
		respondWithError(w, http.StatusBadRequest, "Nama site tidak valid")
		return
	}

	targetDir := filepath.Join("./www", siteName)
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		respondWithError(w, http.StatusNotFound, "Website tidak ditemukan")
		return
	}

	err := os.RemoveAll(targetDir)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Gagal menghapus folder website: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"status":  "sukses",
		"message": fmt.Sprintf("Website '%s' berhasil dihapus!", siteName),
	})
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Clean(f.Name)
		if strings.HasPrefix(fpath, "..") || filepath.IsAbs(fpath) {
			continue // Zip Slip validation
		}

		targetPath := filepath.Join(dest, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(targetPath, f.Mode())
			continue
		}

		if err = os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

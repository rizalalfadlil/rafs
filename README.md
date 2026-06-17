# 🚀 Local Development Server

Proyek ini adalah fondasi server lokal (self-hosted) yang dibangun di dalam lingkungan Docker menggunakan bahasa pemrograman Golang. Server ini dirancang khusus untuk keperluan belajar infrastruktur IT, backend development, dan hosting tugas-tugas proyek perkuliahan secara terisolasi.

## 🏗️ Arsitektur & Teknologi

**Bahasa Inti**: Golang

**Containerization**: Docker & Docker Compose

**Database Induk**: PostgreSQL (berjalan sebagai service terpisah di Docker)

**Web Server**: Golang http.FileServer terhubung dengan Docker Volumes

## 📂 Struktur Direktori

```
rafs/
├── www/                  # (Docker Volume) Tempat menyimpan folder web statis
│   ├── dbui/             # Halaman khusus untuk mengelola database dalam GUI
│   │   └── index.html
│   └── project-lain/
│       └── index.html
├── main.go               # Kode sumber utama server Golang
├── databases/            # Kumpulan kode yang berhubungan dengan database
    ├── database.go
    └── table.go
├── Dockerfile            # Instruksi multi-stage build untuk aplikasi Go
├── docker-compose.yml    # Orkestrasi container web dan database
└── README.md             # Dokumentasi ini

```

## 🌟 Fitur Utama & Progress Roadmap

Proyek ini difokuskan pada tiga pilar utama untuk mendukung perkuliahan. Berikut adalah status pengembangan dari masing-masing fitur:

### 🌐 1. Hosting Web Statis

Fitur untuk menampilkan hasil pengerjaan tugas web (HTML/CSS/JS) secara instan tanpa perlu melakukan konfigurasi server berulang kali.

Daftar Rencana (To-Do):

- [x] Setup sistem routing file statis dasar.

- [x] Integrasi Docker Volume agar tidak perlu rebuild container.

- [ ] Mengintegrasikan CI/CD (GitHub Actions) untuk otomatisasi deployment.

### 🗄️ 2. Pengelola Database Dinamis

Fitur yang bertindak sebagai "Pusat Database". Jika ada tugas baru yang butuh database, server ini bisa otomatis membuatkan database dan user baru yang terisolasi.

Daftar Rencana (To-Do):

- [x] Setup PostgreSQL Master container di docker-compose.yml.

- [x] Pembuatan API POST /api/create-db di Golang.

- [-] Dashboard GUI: Membuat halaman khusus untuk mengelola database dalam GUI **(berhasil dibuat hingga tahap crud database dan tabel, tampilan dan edit untuk field dan record belum diimplementasikan).**

- [ ] Meningkatkan Keamanan dengan menambahkan validasi input keamanan, serta sistem autentikasi API Key atau Basic Auth pada endpoint sensitif agar tidak bisa diakses sembarang orang.

- [ ] Menambahkan fitur untuk menjalankan perintah SQL dalam GUI agar bisa mengelola database lebih mudah.

### 🧩 3. Kumpulan API Service

Tempat mengumpulkan berbagai backend service atau endpoint REST API yang nanti dibangun menggunakan Golang untuk mendukung tugas frontend atau aplikasi mobile.

Daftar Rencana (To-Do):

- [ ] Struktur dasar router Golang.

- [ ] Sistem Autentikasi API: Menambahkan pengamanan API Key atau Basic Auth pada endpoint sensitif agar tidak bisa diakses sembarang orang.

- [ ] Server Monitoring API: Membuat endpoint /api/status yang dapat membaca penggunaan RAM/CPU komputer dan menampilkannya di Dashboard secara real-time.

## 🚀 Cara Menjalankan Server

Pastikan Docker Desktop atau Docker Engine sudah menyala di komputermu.

Buka terminal di direktori proyek ini.

Jalankan perintah untuk merakit dan menyalakan kontainer di latar belakang:

```
docker compose up -d --build
```


Akses halaman welcome/dashboard di: http://localhost:8080/sites/<nama-folder-dashboard>/
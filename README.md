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
├── admin/                # Frontend Panel Admin (Vue 3 + Vite + Tailwind + PrimeVue)
│   ├── dist/             # Build output dari web admin panel (disajikan di /admin/)
│   ├── src/
│   │   ├── components/   # Komponen UI modular (Dialogs, Sidebar, Details)
│   │   ├── composables/  # Logika bisnis/state terpisah (useDatabase.js)
│   │   ├── pages/        # Halaman utama (dashboard.vue, database.vue)
│   │   └── App.vue
│   ├── vite.config.js
│   └── package.json
├── databases/            # Backend package untuk manajemen database, tabel, kolom, & baris
│   ├── database.go       # Operasi CRUD database
│   ├── table.go          # Operasi CRUD tabel
│   ├── column.go         # Operasi CRUD kolom (ADD, DROP, RENAME)
│   ├── row.go            # Operasi CRUD baris (INSERT, UPDATE, DELETE)
│   └── system.go         # Collector metrik sistem real-time (/api/server-info)
├── www/                  # Tempat menyimpan folder web statis (sites/about, sites/hello, dll.)
├── main.go               # Kode sumber utama server Golang
├── Dockerfile            # Instruksi multi-stage build untuk aplikasi Go (Vue + Go runtime)
├── docker-compose.yml    # Orkestrasi container web dan database
└── README.md             # Dokumentasi ini
```

## 🌟 Progress Roadmap Fitur Utama (TODO List)

### 🌐 1. Hosting Web Statis

- [x] Menampilkan halaman untuk setiap folder yang memiliki `index.html` dalam direktori web statis.

- [ ] Fitur upload folder/clone github untuk membuat halaman baru.

- [ ] Integrasi CI/CD (GitHub Actions) untuk otomatisasi deployment.

### 🗄️ 2. Pengelola Database Berbasis GUI

- [x] Setup PostgreSQL Master container di docker-compose.yml.

- [x] Pembuatan API POST /api/create-db di Golang.

- [x] Dashboard GUI: Membuat halaman khusus untuk mengelola database dalam GUI (CRUD database, tabel, kolom, dan baris data).

- [ ] Meningkatkan Keamanan dengan menambahkan validasi input keamanan, serta sistem autentikasi API Key atau Basic Auth pada endpoint sensitif agar tidak bisa diakses sembarang orang.

- [ ] Menambahkan fitur untuk menjalankan perintah SQL dalam GUI agar bisa mengelola database lebih mudah.

### 🌐 3. Kumpulan API Service

Tempat mengumpulkan berbagai backend service atau endpoint REST API yang nanti dibangun menggunakan Golang untuk mendukung tugas frontend atau aplikasi mobile.

Daftar Rencana (To-Do):

- [ ] Struktur dasar router Golang.

- [ ] Sistem Autentikasi API: Menambahkan pengamanan API Key atau Basic Auth pada endpoint sensitif agar tidak bisa diakses sembarang orang.

- [x] Server Monitoring API: Membuat endpoint /api/server-info yang dapat membaca penggunaan RAM/CPU komputer dan menampilkannya di Dashboard secara real-time.

### ☁ 4. Cloud Storage Pribadi

- [ ] Server dapat menyimpan file untuk keperluan apa saja.

- [ ] File-file tertentu dapat diakses melalui API untuk ditampilkan di web.

## 🚀 Cara Menjalankan Server

Pastikan Docker Desktop atau Docker Engine sudah menyala di komputermu.

Buka terminal di direktori proyek ini.

Jalankan perintah untuk merakit dan menyalakan kontainer di latar belakang:

```
docker compose up -d --build
```


Akses halaman welcome/dashboard di: http://localhost:8080/sites/<nama-folder-dashboard>/
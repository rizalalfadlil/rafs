# 📘 Panduan API Server RAFS

Dokumen ini menyediakan spesifikasi lengkap untuk seluruh endpoint REST API yang tersedia pada server **RAFS (Local Development Server)**. Server ini dibangun menggunakan bahasa pemrograman Go (Golang) dan menyediakan antarmuka API untuk manajemen database PostgreSQL, hosting website statis, pemantauan performa sistem, serta cloud storage pribadi.

---

## 📌 Informasi Umum
* **Base URL**: `http://localhost:8080` (atau port host yang ditentukan)
* **Format Response**: Semua response REST API mengembalikan format data **JSON** kecuali untuk endpoint download.
* **Format Request**: Seluruh request bertipe `POST`/`PUT` membutuhkan header `Content-Type: application/json` kecuali untuk endpoint yang memerlukan unggahan berkas (`multipart/form-data`).

---

## 🗺️ Ringkasan Endpoint

| Kategori | Metode | Endpoint | Deskripsi |
| :--- | :--- | :--- | :--- |
| **Database** | `POST` | `/api/create-db` | Membuat database dan user PostgreSQL (Compat) |
| | `GET` | `/api/databases` | List database buatan pengguna |
| | `POST` | `/api/databases` | Membuat database dan user PostgreSQL baru |
| | `PUT` | `/api/databases` | Mengubah nama (rename) database |
| | `DELETE`| `/api/databases` | Menghapus database dan memutuskan koneksinya |
| **Tabel** | `GET` | `/api/tables` | List tabel dalam suatu database |
| | `POST` | `/api/tables` | Membuat tabel baru beserta kolomnya |
| | `PUT` | `/api/tables` | Mengubah nama (rename) tabel |
| | `DELETE`| `/api/tables` | Menghapus tabel dari database |
| **Kolom** | `GET` | `/api/columns` | List kolom dan skema dari suatu tabel |
| | `POST` | `/api/columns` | Menambahkan kolom baru ke tabel |
| | `PUT` | `/api/columns` | Mengubah nama/tipe data kolom |
| | `DELETE`| `/api/columns` | Menghapus kolom dari tabel |
| **Baris (Data)**| `GET` | `/api/rows` | Mengambil seluruh baris data tabel |
| | `POST` | `/api/rows` | Menyisipkan baris baru (Insert) |
| | `PUT` | `/api/rows` | Memperbarui baris berdasarkan kriteria (Update) |
| | `DELETE`| `/api/rows` | Menghapus baris berdasarkan kriteria (Delete) |
| **Sistem** | `GET` | `/api/server-info` | Informasi metrik server real-time |
| **Situs Statis**| `GET` | `/api/sites` | List website statis di folder `./www` |
| | `POST` | `/api/sites/clone` | Clone situs statis dari repositori Git publik |
| | `POST` | `/api/sites/upload`| Unggah dan ekstrak file zip situs statis |
| | `DELETE`| `/api/sites` | Menghapus folder situs statis |
| **Storage** | `GET` | `/api/storage` | List file/folder & status kuota cloud storage |
| | `POST` | `/api/storage/folder` | Membuat folder baru di storage |
| | `POST` | `/api/storage/upload` | Unggah file (multiple) ke storage |
| | `POST` | `/api/storage/delete` | Hapus file/folder secara rekursif |
| | `GET` | `/api/storage/download`| Unduh/lihat file secara langsung |
| | `POST` | `/api/storage/public` | Publikasikan file/folder ke direktori publik |

---

## 🗄️ 1. Manajemen Database & Pengguna

### 1.1 List Database
Mengambil daftar database yang dibuat oleh pengguna (tidak termasuk sistem internal PostgreSQL).

* **Method**: `GET`
* **URL**: `/api/databases`
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "databases": ["my_db", "test_project"]
  }
  ```

### 1.2 Membuat Database Baru
Membuat user PostgreSQL baru sekaligus database baru dengan kepemilikan (owner) di bawah user tersebut.

* **Method**: `POST`
* **URL**: `/api/databases` (Atau `/api/create-db` untuk kompatibilitas ke belakang)
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "username": "project_user",
    "password": "securepassword123"
  }
  ```
  > [!NOTE]
  > Karakter `db_name` dan `username` divalidasi menggunakan ekspresi reguler `^[a-zA-Z_][a-zA-Z0-9_]*$` untuk mencegah SQL injection.
* **Response (210 Created)**:
  ```json
  {
    "status": "sukses",
    "message": "Database 'project_db' dan User 'project_user' berhasil dibuat!"
  }
  ```

### 1.3 Mengubah Nama Database
Mengubah identitas nama database yang ada.

* **Method**: `PUT`
* **URL**: `/api/databases`
* **Request Body (JSON)**:
  ```json
  {
    "old_name": "project_db",
    "new_name": "new_project_db"
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "Database 'project_db' berhasil diubah nama menjadi 'new_project_db'!"
  }
  ```

### 1.4 Menghapus Database
Menghapus database beserta koneksi aktif yang masih melekat pada database tersebut.

* **Method**: `DELETE`
* **URL**: `/api/databases`
* **Query Parameters**:
  * `db_name`: `new_project_db` *(Wajib)*
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "Database 'new_project_db' berhasil dihapus!"
  }
  ```

---

## 📋 2. Manajemen Tabel

### 2.1 List Tabel
Mengambil semua tabel yang berada di skema `public` dalam database tertentu.

* **Method**: `GET`
* **URL**: `/api/tables`
* **Query Parameters**:
  * `db_name`: `project_db` *(Wajib)*
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "tables": ["users", "products"]
  }
  ```

### 2.2 Membuat Tabel Baru
Membuat tabel baru beserta spesifikasi kolom awalnya.

* **Method**: `POST`
* **URL**: `/api/tables`
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "table_name": "users",
    "columns": [
      { "name": "id", "type": "SERIAL PRIMARY KEY" },
      { "name": "username", "type": "VARCHAR(100) NOT NULL" },
      { "name": "email", "type": "VARCHAR(255)" }
    ]
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "status": "sukses",
    "message": "Tabel 'users' berhasil dibuat di database 'project_db'!"
  }
  ```

### 2.3 Mengubah Nama Tabel
Mengubah nama tabel yang ada dalam database.

* **Method**: `PUT`
* **URL**: `/api/tables`
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "old_name": "users",
    "new_name": "accounts"
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "Tabel 'users' berhasil diubah nama menjadi 'accounts' di database 'project_db'!"
  }
  ```

### 2.4 Menghapus Tabel
Menghapus tabel tertentu secara permanen.

* **Method**: `DELETE`
* **URL**: `/api/tables`
* **Query Parameters**:
  * `db_name`: `project_db` *(Wajib)*
  * `table_name`: `accounts` *(Wajib)*
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "Tabel 'accounts' berhasil dihapus dari database 'project_db'!"
  }
  ```

---

## 📑 3. Skema Kolom Tabel

### 3.1 List Kolom
Mengambil metadata skema kolom dari tabel yang dipilih.

* **Method**: `GET`
* **URL**: `/api/columns`
* **Query Parameters**:
  * `db_name`: `project_db` *(Wajib)*
  * `table_name`: `users` *(Wajib)*
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "columns": [
      {
        "name": "id",
        "type": "integer",
        "nullable": "NO",
        "default": "nextval('users_id_seq'::regclass)"
      },
      {
        "name": "username",
        "type": "character varying",
        "nullable": "NO",
        "default": ""
      }
    ]
  }
  ```

### 3.2 Menambahkan Kolom Baru
Menyisipkan kolom baru ke dalam tabel yang sudah ada.

* **Method**: `POST`
* **URL**: `/api/columns`
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "table_name": "users",
    "column_name": "created_at",
    "column_type": "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "status": "sukses",
    "message": "Kolom 'created_at' berhasil ditambahkan ke tabel 'users'!"
  }
  ```

### 3.3 Mengubah Kolom (Rename & Ganti Tipe)
Mengubah nama kolom, tipe data kolom, atau keduanya.

* **Method**: `PUT`
* **URL**: `/api/columns`
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "table_name": "users",
    "old_name": "username",
    "new_name": "user_name",
    "column_type": "VARCHAR(150)"
  }
  ```
  > [!TIP]
  > Properti `column_type` bersifat opsional. Jika hanya ingin me-rename, kosongkan properti tersebut.
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "Kolom 'username' berhasil diubah!"
  }
  ```

### 3.4 Menghapus Kolom
Menghapus kolom tertentu dari skema tabel.

* **Method**: `DELETE`
* **URL**: `/api/columns`
* **Query Parameters**:
  * `db_name`: `project_db` *(Wajib)*
  * `table_name`: `users` *(Wajib)*
  * `column_name`: `created_at` *(Wajib)*
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "Kolom 'created_at' berhasil dihapus dari tabel 'users'!"
  }
  ```

---

## 📊 4. Manipulasi Data Baris (Rows CRUD)

### 4.1 List Baris Data
Mengambil semua baris data yang ada pada suatu tabel.

* **Method**: `GET`
* **URL**: `/api/rows`
* **Query Parameters**:
  * `db_name`: `project_db` *(Wajib)*
  * `table_name`: `users` *(Wajib)*
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "rows": [
      {
        "id": 1,
        "username": "rizal",
        "email": "rizal@example.com"
      }
    ]
  }
  ```

### 4.2 Menyisipkan Baris Data Baru (Insert)
Menambahkan data record baru ke dalam kolom-kolom tabel.

* **Method**: `POST`
* **URL**: `/api/rows`
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "table_name": "users",
    "row": {
      "username": "alfadlil",
      "email": "alfadlil@example.com"
    }
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "status": "sukses",
    "message": "Baris baru berhasil ditambahkan!"
  }
  ```

### 4.3 Memperbarui Baris Data (Update)
Memperbarui nilai kolom untuk baris data yang memenuhi filter kriteria (`where`).

* **Method**: `PUT`
* **URL**: `/api/rows`
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "table_name": "users",
    "where": {
      "id": 1
    },
    "row": {
      "username": "rizal_new",
      "email": "rizal_new@example.com"
    }
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "Baris berhasil diperbarui!"
  }
  ```

### 4.4 Menghapus Baris Data (Delete)
Menghapus baris data yang memenuhi filter kriteria (`where`).

* **Method**: `DELETE`
* **URL**: `/api/rows`
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "table_name": "users",
    "where": {
      "id": 2
    }
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "Baris berhasil dihapus!"
  }
  ```

---

## 📈 5. Monitoring Kinerja Sistem

### 5.1 Real-Time Server Info
Mengambil metrik performa komputer server secara langsung seperti Uptime, Penggunaan Memori, Penggunaan CPU, Sistem Operasi, Versi Compiler Go, dan Versi Database PostgreSQL.

* **Method**: `GET`
* **URL**: `/api/server-info`
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "info": {
      "uptime": "1 Hari 5 Jam 23 Menit 12 Detik",
      "memory_used": "384 MB",
      "memory_total": "16384 MB",
      "memory_percentage": "2.3%",
      "cpu_usage": "1.2%",
      "os": "Ubuntu 22.04.3 LTS",
      "go_version": "go1.21.5",
      "postgres_version": "PostgreSQL 15.2"
    }
  }
  ```

---

## 🌐 6. Hosting Web Statis (Static Website)

Seluruh situs yang dikelola akan disimpan di direktori `./www/<site_name>` dan disajikan secara publik melalui alamat URL `http://localhost:8080/sites/<site_name>/`.

### 6.1 List Website Statis
Menampilkan daftar website statis yang saat ini di-hosting beserta status keaktifan (adanya file `index.html`).

* **Method**: `GET`
* **URL**: `/api/sites`
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "sites": [
      {
        "name": "my-portfolio",
        "active": true
      },
      {
        "name": "landing-page",
        "active": false
      }
    ]
  }
  ```

### 6.2 Clone Website dari Git
Melakukan `git clone` dari repositori publik untuk dijadikan situs statis baru.

* **Method**: `POST`
* **URL**: `/api/sites/clone`
* **Request Body (JSON)**:
  ```json
  {
    "repo_url": "https://github.com/username/repo-name.git",
    "site_name": "portfolio-github"
  }
  ```
  > [!WARNING]
  > Direktori `.git` akan dihapus setelah proses klon selesai untuk menjaga kebersihan penyimpanan server.
* **Response (201 Created)**:
  ```json
  {
    "status": "sukses",
    "message": "Repository berhasil diclone ke website 'portfolio-github'!",
    "active": true
  }
  ```

### 6.3 Upload Website via File ZIP
Menerima berkas arsip `.zip` dan mengekstraknya sebagai website statis.

* **Method**: `POST`
* **URL**: `/api/sites/upload`
* **Content-Type**: `multipart/form-data`
* **Payload (Form-Data)**:
  * `site_name`: `my-app` *(Text)*
  * `file`: `arsip-situs.zip` *(File, maks. 50 MB)*
* **Response (210 Created)**:
  ```json
  {
    "status": "sukses",
    "message": "Website 'my-app' berhasil diunggah!",
    "active": true
  }
  ```

### 6.4 Menghapus Website Statis
Menghapus website statis beserta semua berkasnya secara permanen dari server.

* **Method**: `DELETE`
* **URL**: `/api/sites`
* **Query Parameters**:
  * `site_name`: `landing-page` *(Wajib)*
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "Website 'landing-page' berhasil dihapus!"
  }
  ```

---

## ☁️ 7. Cloud Storage Pribadi

Layanan cloud storage pribadi yang mendukung manipulasi berkas, pembuatan folder, pengunggahan, pemantauan batas kuota penyimpanan, dan opsi pemublikasian tautan unduhan.

### 7.1 List Berkas & Status Kuota
Menampilkan isi berkas dan subfolder dalam direktori saat ini, beserta data penggunaan penyimpanan (Kuota Storage dikunci pada batas maksimal **1 GB**).

* **Method**: `GET`
* **URL**: `/api/storage`
* **Query Parameters**:
  * `path`: `tugas/kuliah` *(Opsional, kosongkan untuk menampilkan root storage)*
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "current_path": "tugas/kuliah",
    "contents": [
      {
        "name": "laporan.pdf",
        "type": "file",
        "size": 2548291,
        "mod_time": "2026-06-21 09:21:40"
      },
      {
        "name": "referensi_gambar",
        "type": "folder",
        "size": 0,
        "mod_time": "2026-06-21 08:15:00"
      }
    ],
    "space": {
      "used": 45182900,
      "total": 1073741824
    }
  }
  ```

### 7.2 Membuat Folder Baru
Membuat subdirektori baru dalam lokasi penyimpanan cloud.

* **Method**: `POST`
* **URL**: `/api/storage/folder`
* **Request Body (JSON)**:
  ```json
  {
    "path": "tugas/kuliah",
    "name": "sumber_data"
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "status": "sukses",
    "message": "Folder 'sumber_data' berhasil dibuat!"
  }
  ```

### 7.3 Mengunggah File (Multiple Upload)
Mengunggah satu atau lebih berkas ke lokasi folder penyimpanan tertentu.

* **Method**: `POST`
* **URL**: `/api/storage/upload`
* **Content-Type**: `multipart/form-data`
* **Payload (Form-Data)**:
  * `path`: `tugas/kuliah/sumber_data` *(Text)*
  * `files`: Berkas data 1 *(File, batas maks. total upload 100 MB)*
  * `files`: Berkas data 2 *(File, nama field harus 'files' untuk multi-upload)*
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "2 berkas berhasil diunggah!"
  }
  ```

### 7.4 Mengunduh / Preview Berkas
Mendapatkan konten berkas secara langsung untuk proses pengunduhan atau pratinjau (preview).

* **Method**: `GET`
* **URL**: `/api/storage/download`
* **Query Parameters**:
  * `path`: `tugas/kuliah/laporan.pdf` *(Wajib)*
* **Response (200 OK)**:
  * Stream Binary Berkas (menggunakan `http.ServeFile` di Go) dengan header Content-Type otomatis sesuai tipe berkas.

### 7.5 Menghapus Berkas / Folder
Menghapus satu atau beberapa berkas/folder secara rekursif.

* **Method**: `POST`
* **URL**: `/api/storage/delete`
* **Request Body (JSON)**:
  ```json
  {
    "paths": [
      "tugas/kuliah/laporan.pdf",
      "tugas/kuliah/sumber_data"
    ]
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "2 file/folder berhasil dihapus!"
  }
  ```

### 7.6 Publikasikan Berkas ke Folder Publik
Menyalin berkas atau folder secara rekursif dari direktori cloud privat ke direktori `/public` dan menghasilkan tautan URL publik yang dapat diakses dari mana saja secara langsung.

* **Method**: `POST`
* **URL**: `/api/storage/public`
* **Request Body (JSON)**:
  ```json
  {
    "paths": [
      "tugas/kuliah/laporan.pdf"
    ]
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "1 file/folder berhasil dipublikasikan!",
    "links": [
      "/public/tugas/kuliah/laporan.pdf"
    ]
  }
  ```
  > [!TIP]
  > Tautan hasil publikasi di atas dapat diakses langsung menggunakan URL host: `http://localhost:8080/public/tugas/kuliah/laporan.pdf`

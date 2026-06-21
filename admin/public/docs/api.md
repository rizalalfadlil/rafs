# 📘 Panduan API Server RAFS

Dokumen ini menyediakan spesifikasi lengkap untuk seluruh endpoint REST API yang tersedia pada server **RAFS (Local Development Server)**. Server ini dibangun menggunakan bahasa pemrograman Go (Golang) dan menyediakan antarmuka API untuk manajemen database PostgreSQL, hosting website statis, pemantauan performa sistem, serta cloud storage pribadi.

---

## 📌 Informasi Umum
* **Base URL**: `http://localhost:8080` (atau port host yang ditentukan)
* **Format Response**: Semua response REST API mengembalikan format data **JSON** kecuali untuk endpoint download.
* **Format Request**: Seluruh request bertipe `POST`/`PUT` membutuhkan header `Content-Type: application/json` kecuali untuk endpoint yang memerlukan unggahan berkas (`multipart/form-data`).

---

## 🔒 Autentikasi Koneksi Database

> [!IMPORTANT]
> Untuk meningkatkan keamanan, seluruh operasi akses data dan struktur tabel mewajibkan pengiriman kredensial pemilik database (Owner) yang sah menggunakan custom request headers:
> * **`X-Database-User`**: Username PostgreSQL pemilik database.
> * **`X-Database-Password`**: Password PostgreSQL pemilik database.
>
> Header ini wajib disertakan untuk semua pemanggilan endpoint berikut:
> * **Database**: `PUT /api/databases` (Rename), `DELETE /api/databases` (Drop)
> * **Tabel**: `GET /api/tables` (List), `POST /api/tables` (Create), `PUT /api/tables` (Rename), `DELETE /api/tables` (Drop)
> * **Kolom**: `GET /api/columns`, `POST /api/columns`, `PUT /api/columns`, `DELETE /api/columns`
> * **Baris (Rows)**: `GET /api/rows`, `POST /api/rows`, `PUT /api/rows`, `DELETE /api/rows`
> * **Raw SQL Command**: `POST /api/query` (SQL Executor)

---

## 🗺️ Ringkasan Endpoint

| Kategori | Metode | Endpoint | Deskripsi | Autentikasi |
| :--- | :--- | :--- | :--- | :--- |
| **Database** | `POST` | `/api/create-db` | Membuat database dan user PostgreSQL (Compat) | Superadmin (Internal) |
| | `GET` | `/api/databases` | List database buatan pengguna | Superadmin (Internal) |
| | `POST` | `/api/databases` | Membuat database dan user PostgreSQL baru | Superadmin (Internal) |
| | `PUT` | `/api/databases` | Mengubah nama (rename) database | **Header Owner** |
| | `DELETE`| `/api/databases` | Menghapus database dan memutuskan koneksinya | **Header Owner** |
| **Tabel** | `GET` | `/api/tables` | List tabel dalam suatu database | **Header Owner** |
| | `POST` | `/api/tables` | Membuat tabel baru beserta kolomnya | **Header Owner** |
| | `PUT` | `/api/tables` | Mengubah nama (rename) tabel | **Header Owner** |
| | `DELETE`| `/api/tables` | Menghapus tabel dari database | **Header Owner** |
| **Kolom** | `GET` | `/api/columns` | List kolom dan skema dari suatu tabel | **Header Owner** |
| | `POST` | `/api/columns` | Menambahkan kolom baru ke tabel | **Header Owner** |
| | `PUT` | `/api/columns` | Mengubah nama/tipe data kolom | **Header Owner** |
| | `DELETE`| `/api/columns` | Menghapus kolom dari tabel | **Header Owner** |
| **Baris (Data)**| `GET` | `/api/rows` | Mengambil seluruh baris data tabel | **Header Owner** |
| | `POST` | `/api/rows` | Menyisipkan baris baru (Insert) | **Header Owner** |
| | `PUT` | `/api/rows` | Memperbarui baris berdasarkan kriteria (Update) | **Header Owner** |
| | `DELETE`| `/api/rows` | Menghapus baris berdasarkan kriteria (Delete) | **Header Owner** |
| | `POST` | `/api/query` | Mengeksekusi perintah SQL mentah (SELECT/DDL/DML) | **Header Owner** |
| **Sistem** | `GET` | `/api/server-info` | Informasi metrik server real-time | None |
| **Situs Statis**| `GET` | `/api/sites` | List website statis di folder `./www` | None |
| | `POST` | `/api/sites/clone` | Clone situs statis dari repositori Git publik | None |
| | `POST` | `/api/sites/upload`| Unggah dan ekstrak file zip situs statis | None |
| | `DELETE`| `/api/sites` | Menghapus folder situs statis | None |
| **Storage** | `GET` | `/api/storage` | List file/folder & status kuota cloud storage | None |
| | `POST` | `/api/storage/folder` | Membuat folder baru di storage | None |
| | `POST` | `/api/storage/upload` | Unggah file (multiple) ke storage | None |
| | `POST` | `/api/storage/delete` | Hapus file/folder secara rekursif | None |
| | `GET` | `/api/storage/download`| Unduh/lihat file secara langsung | None |
| | `POST` | `/api/storage/public` | Publikasikan file/folder ke direktori publik | None |

---

## 🗄️ 1. Manajemen Database & Pengguna

### 1.1 List Database
Mengambil daftar database yang dibuat oleh pengguna.

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
* **URL**: `/api/databases`
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "username": "project_user",
    "password": "securepassword123"
  }
  ```
* **Response (201 Created)**:
  ```json
  {
    "status": "sukses",
    "message": "Database 'project_db' dan User 'project_user' berhasil dibuat!"
  }
  ```

---

## 📋 2. Manajemen Tabel

### 2.1 List Tabel
Mengambil semua tabel yang berada di skema `public` dalam database tertentu.

* **Method**: `GET`
* **URL**: `/api/tables`
* **Headers**:
  * `X-Database-User`: `project_user`
  * `X-Database-Password`: `securepassword123`
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
* **Headers**:
  * `X-Database-User`: `project_user`
  * `X-Database-Password`: `securepassword123`
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "table_name": "users",
    "columns": [
      { "name": "id", "type": "SERIAL PRIMARY KEY" },
      { "name": "username", "type": "VARCHAR(100) NOT NULL" }
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

---

## 💻 3. Eksekusi SQL Custom

Menerima input query SQL mentah dan mengeksekusinya di bawah koneksi pengguna yang sesuai.

* **Method**: `POST`
* **URL**: `/api/query`
* **Headers**:
  * `X-Database-User`: `project_user`
  * `X-Database-Password`: `securepassword123`
* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "query": "SELECT * FROM users ORDER BY id DESC;"
  }
  ```

### 3.1 Respons Sukses (SELECT Query)
Kembali dengan status 200 OK beserta list nama kolom (`columns`) dan data baris (`rows`).

* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "type": "select",
    "columns": ["id", "username"],
    "rows": [
      { "id": 1, "username": "admin" },
      { "id": 2, "username": "guest" }
    ]
  }
  ```

### 3.2 Respons Sukses (DDL/DML Exec)
Kembali dengan status 200 OK beserta jumlah baris data yang terpengaruh (`rows_affected`).

* **Request Body (JSON)**:
  ```json
  {
    "db_name": "project_db",
    "query": "UPDATE users SET username = 'super_admin' WHERE id = 1;"
  }
  ```
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "type": "exec",
    "rows_affected": 1,
    "message": "Perintah berhasil dieksekusi."
  }
  ```

### 3.3 Respons Gagal (Error SQL)
Jika sintaks SQL salah, server mengembalikan status 400 Bad Request beserta pesan error dari engine PostgreSQL.

* **Response (400 Bad Request)**:
  ```json
  {
    "status": "error",
    "message": "pq: relation \"user_profile\" does not exist"
  }
  ```

---

## 🌐 4. Hosting Web Statis (Static Website)

### 4.1 List Website Statis
Menampilkan daftar website statis yang saat ini di-hosting.

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
      }
    ]
  }
  ```

### 4.2 Clone Website dari Git
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
* **Response (201 Created)**:
  ```json
  {
    "status": "sukses",
    "message": "Repository berhasil diclone ke website 'portfolio-github'!",
    "active": true
  }
  ```

---

## ☁️ 5. Cloud Storage Pribadi

### 5.1 List Berkas & Status Kuota
Menampilkan isi berkas dan subfolder dalam direktori saat ini.

* **Method**: `GET`
* **URL**: `/api/storage`
* **Query Parameters**:
  * `path`: `tugas/kuliah` *(Opsional)*
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
      }
    ],
    "space": {
      "used": 45182900,
      "total": 1073741824
    }
  }
  ```

### 5.2 Mengunggah File (Multiple Upload)
Mengunggah satu atau lebih berkas ke lokasi folder penyimpanan tertentu.

* **Method**: `POST`
* **URL**: `/api/storage/upload`
* **Content-Type**: `multipart/form-data`
* **Payload (Form-Data)**:
  * `path`: `tugas/kuliah` *(Text)*
  * `files`: Berkas data *(File)*
* **Response (200 OK)**:
  ```json
  {
    "status": "sukses",
    "message": "1 berkas berhasil diunggah!"
  }
  ```

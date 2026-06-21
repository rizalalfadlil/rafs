# 📱 Panduan Setup & Menjalankan Aplikasi di Android via PRoot (Termux)

Panduan ini menjelaskan cara melepas ketergantungan Docker dan menjalankan proyek ini secara lokal langsung di perangkat Android menggunakan **PRoot Distro** di **Termux**. Kita akan menggunakan **Alpine Linux** di dalam PRoot karena sangat ringan dan memiliki lingkungan yang sama dengan kontainer Docker asli proyek ini.

---

## 🛠️ Prasyarat (Prerequisites)

1. **Termux**: Unduh dan pasang aplikasi Termux terbaru dari [F-Droid](https://f-droid.org/en/packages/com.termux/) (jangan gunakan versi Google Play Store karena sudah usang/tidak didevelop lagi).
2. **Koneksi Internet**: Dibutuhkan untuk mengunduh paket Linux dan dependensi proyek.

---

## 🚀 Langkah-langkah Setup

### Langkah 1: Pasang & Jalankan PRoot Distro (Alpine Linux)

Buka Termux, lalu jalankan perintah berikut untuk memperbarui paket dan menginstal `proot-distro` dengan Alpine Linux:

```bash
# 1. Update package manager Termux
pkg update -y && pkg upgrade -y

# 2. Install proot-distro
pkg install proot-distro -y

# 3. Install Alpine Linux di dalam PRoot
proot-distro install alpine

# 4. Masuk (login) ke lingkungan Alpine Linux
proot-distro login alpine
```

> [!TIP]
> Setelah berhasil login, prompt terminal Anda akan berubah (misalnya menjadi `localhost:~#`). Anda sekarang berada di dalam sistem Linux Alpine terisolasi di Android Anda!

---

### Langkah 2: Install Dependensi di Alpine Linux

Di dalam lingkungan Alpine Linux, jalankan perintah berikut untuk menginstal seluruh alat yang dibutuhkan:

```bash
apk update
apk add go nodejs npm postgresql postgresql-client git bash openrc
```

---

### Langkah 3: Setup & Jalankan PostgreSQL

Karena PRoot tidak mendukung systemd atau daemon initialization penuh secara langsung, kita perlu menginisialisasi database secara manual dan menjalankannya di latar belakang menggunakan utilitas `pg_ctl`.

```bash
# 1. Buat folder penyimpanan data PostgreSQL
mkdir -p /var/lib/postgresql/data
chown -R postgres:postgres /var/lib/postgresql/data

# 2. Pindah ke user 'postgres' dan lakukan initdb
su - postgres -c "initdb -D /var/lib/postgresql/data"

# 3. Jalankan server PostgreSQL di latar belakang
su - postgres -c "pg_ctl -D /var/lib/postgresql/data -l /var/lib/postgresql/logfile start"
```

Untuk memastikan PostgreSQL sudah berjalan, Anda bisa mengecek statusnya dengan:
```bash
su - postgres -c "pg_ctl -D /var/lib/postgresql/data status"
```

---

### Langkah 4: Buat Kredensial Superadmin Database

Aplikasi memerlukan user superadmin dengan password khusus agar backend dapat membuat database dan tabel baru secara dinamis melalui UI. Jalankan perintah berikut untuk masuk ke shell PostgreSQL:

```bash
su - postgres -c "psql"
```

Di dalam shell `psql` (prompt `postgres=#`), jalankan perintah SQL berikut:

```sql
-- 1. Buat user superadmin dengan password default proyek
CREATE ROLE superadmin WITH LOGIN PASSWORD 'supersecret123' SUPERUSER;

-- 2. Buat database default 'postgres' (jika belum ada)
CREATE DATABASE postgres OWNER superadmin;

-- 3. Keluar dari psql
\q
```

---

### Langkah 5: Clone Repositori Proyek & Pindah ke Branch Android

Kembali ke folder home user root di Alpine Linux (misal `/root`), clone repositori proyek Anda dan berpindah ke branch `feature/proot-android` yang mendukung variabel lingkungan:

```bash
cd ~
git clone https://github.com/rizalalfadlil/rafs.git
cd rafs
git checkout feature/proot-android
```

---

### Langkah 6: Build Frontend Panel Admin (Vue App)

Sebelum server backend dijalankan, kita perlu membuild aset statis frontend panel admin agar disajikan dengan benar oleh server utama:

```bash
# Masuk ke direktori frontend admin
cd admin

# Install dependensi frontend
npm install

# Build frontend untuk produksi
npm run build

# Kembali ke root direktori proyek
cd ..
```

---

### Langkah 7: Jalankan Server Backend Golang dengan Variabel Lingkungan

Sekarang kita akan menjalankan server Golang. Karena PostgreSQL berjalan di host lokal (di dalam PRoot yang sama), kita harus menyetel variabel lingkungan `DB_HOST` ke `127.0.0.1` (bukan `db` seperti di Docker):

```bash
# Eksport variabel lingkungan untuk koneksi DB
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_SUPERUSER=superadmin
export DB_SUPERPASSWORD=supersecret123
export DB_SUPERDB=postgres
export DB_SSLMODE=disable

# Jalankan aplikasi Golang
go run main.go
```

Output terminal akan menampilkan:
```
Server berjalan di port :8080...
```

---

## 🌐 Cara Mengakses Aplikasi

### A. Dari Browser di Perangkat Android yang Sama
Buka browser favorit Anda di HP Android, dan akses URL berikut:
- Dashboard/Admin GUI: [http://localhost:8080/admin/](http://localhost:8080/admin/)
- Endpoint API Utama: [http://localhost:8080/](http://localhost:8080/)

### B. Dari Laptop/PC di Jaringan Wi-Fi yang Sama
Jika Anda ingin mengakses server ini dari PC/Laptop Anda:
1. Cari tahu IP Lokal HP Android Anda. Di Termux (buka sesi terminal baru di luar PRoot), jalankan:
   ```bash
   ifconfig
   ```
   Cari alamat IPv4 di bawah interface `wlan0` (misal: `192.168.1.50`).
2. Di browser PC/Laptop Anda, akses:
   - [http://192.168.1.50:8080/admin/](http://192.168.1.50:8080/admin/)

---

## 🛑 Cara Menghentikan Server Database & Aplikasi

Jika Anda ingin mematikan server database setelah selesai digunakan:

```bash
# Masuk kembali ke PRoot jika sudah keluar
proot-distro login alpine

# Hentikan PostgreSQL
su - postgres -c "pg_ctl -D /var/lib/postgresql/data stop"
```

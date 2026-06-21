# 🗄️ Panduan Manajemen Database

Layanan Database pada server **RAFS** dirancang untuk memudahkan Anda mengelola basis data relasional PostgreSQL melalui antarmuka GUI yang intuitif. Sistem ini berjalan di atas container Docker terisolasi dan mendukung operasi penuh CRUD untuk database, tabel, kolom, hingga data baris.

---

## 🚀 1. Konsep Dasar & Arsitektur

* **Database Engine**: PostgreSQL 15+.
* **Akses Master**: Menggunakan user `superadmin` dengan hak akses penuh untuk melakukan operasi administratif (membuat database, user baru, serta mengelola hak akses).
* **Isolasi Database**: Setiap database baru yang dibuat melalui GUI akan memiliki **User khusus** sebagai pemiliknya (owner). Hal ini bertujuan untuk mengisolasi data antar proyek agar lebih aman.

---

## ➕ 2. Membuat Database Baru

Untuk membuat database baru:
1. Buka halaman **Database Management** di menu panel admin.
2. Klik tombol **Buat Database Baru**.
3. Isi formulir dengan data berikut:
   * **Nama Database**: Identifikasi database (contoh: `toko_online`).
   * **Username Owner**: User yang akan bertindak sebagai pemilik database (contoh: `toko_admin`).
   * **Password**: Kata sandi untuk user owner tersebut.
4. Klik **Simpan**. Sistem akan mengeksekusi pembuatan user dan database baru secara otomatis.

> [!IMPORTANT]
> **Aturan Penamaan (Validasi Input):**
> * Nama database dan username hanya boleh diawali dengan huruf atau underscore (`_`) dan diikuti oleh huruf, angka, atau underscore.
> * Ekspresi Reguler yang digunakan: `^[a-zA-Z_][a-zA-Z0-9_]*$`.
> * Spasi dan karakter khusus seperti `-`, `@`, `$`, `#` **dilarang** demi menghindari celah keamanan SQL Injection.

---

## 🔑 3. Akses & Autentikasi Database

Setiap database yang di-hosting terlindungi secara ketat. Pengguna wajib memasukkan kredensial pemilik database untuk mengakses isinya:

1. **Prompt Login**: Saat Anda memilih database dari menu sidebar, dialog **Akses Database** akan muncul.
2. **Kredensial**: Masukkan **Username PostgreSQL** dan **Password** yang valid sesuai dengan pemilik database tersebut.
3. **Penyimpanan Cache**: Untuk kenyamanan Anda, kredensial yang berhasil diverifikasi akan disimpan secara lokal di dalam cache browser Anda (`localStorage`). Anda tidak perlu memasukkan kata sandi lagi saat memuat ulang halaman ini.
4. **Keluar Database**: Anda dapat memutus sesi koneksi dan menghapus cache kata sandi kapan saja dengan mengklik tombol **Keluar (Sign Out)** di sebelah kanan status koneksi aktif pada footer.

---

## 📊 4. Mengelola Tabel & Skema

Setiap database terdiri dari tabel-tabel tempat menyimpan data. Melalui panel admin, Anda dapat melakukan operasi berikut:

### 4.1 Membuat Tabel Baru
Saat membuat tabel, Anda wajib mendefinisikan kolom awal:
1. Pilih database aktif dari daftar database di sebelah kiri.
2. Klik **Buat Tabel Baru** pada panel tabel.
3. Masukkan **Nama Tabel** dan tambahkan kolom dengan menentukan **Nama Kolom** serta **Tipe Data**.
   * *Contoh kolom 1*: Nama `id`, Tipe data `SERIAL PRIMARY KEY` (otomatis bertambah).
   * *Contoh kolom 2*: Nama `nama_produk`, Tipe data `VARCHAR(100) NOT NULL`.
   * *Contoh kolom 3*: Nama `harga`, Tipe data `INTEGER`.

### 4.2 Tipe Data yang Didukung
Anda dapat menggunakan tipe data standar PostgreSQL, seperti:
* `SERIAL` / `BIGSERIAL`: Bilangan bulat auto-increment (cocok untuk Primary Key).
* `INTEGER` / `BIGINT`: Bilangan bulat.
* `VARCHAR(panjang)` / `TEXT`: Data teks/karakter string.
* `BOOLEAN`: Nilai true atau false.
* `TIMESTAMP` / `DATE`: Data waktu dan tanggal.

---

## 📑 5. Mengelola Kolom (Skema Dinamis)

Anda dapat mengubah struktur tabel kapan saja tanpa perlu menulis query SQL secara manual:
* **Tambah Kolom**: Menambahkan field baru ke tabel yang sudah ada dengan tipe data yang ditentukan.
* **Ubah Kolom**: Mengubah nama kolom (`Rename`) atau mengubah tipe datanya (`Alter Type`).
* **Hapus Kolom**: Menghapus kolom beserta seluruh data yang tersimpan di dalamnya (`Drop Column`).

> [!WARNING]
> Menghapus kolom bersifat **permanen** dan akan melenyapkan seluruh data yang ada di kolom tersebut. Harap lakukan dengan hati-hati.

---

## 🖊️ 6. Manipulasi Data Baris (Rows CRUD)

Setelah tabel dan kolom siap, Anda dapat langsung mengelola data di dalamnya:
* **Insert Row (Tambah Data)**: Mengisi nilai untuk kolom-kolom yang tersedia di tabel melalui form input dinamis.
* **Update Row (Edit Data)**: Memperbarui nilai pada baris tertentu. Sistem secara otomatis menyusun kriteria pencarian (`WHERE`) berdasarkan kolom pengenal unik untuk memastikan baris yang tepat diperbarui.
* **Delete Row (Hapus Data)**: Menghapus baris data terpilih dari tabel.

---

## 💻 7. Terminal SQL Query

Bagi pengguna tingkat lanjut, RAFS menyediakan terminal kueri SQL terintegrasi di bagian bawah halaman pengelola database:

1. **Menjalankan Query**: Masukkan perintah SQL kueri mentah di area input teks. Klik tombol **Run** atau gunakan pintasan keyboard **`Ctrl + Enter`** untuk mengeksekusi perintah.
2. **Keluaran SELECT**: Hasil eksekusi query bertipe `SELECT` akan dirender dalam format tabel monospaced yang interaktif di bawah editor kueri.
3. **Keluaran Exec (DDL/DML)**: Perintah penulisan data seperti `INSERT`, `UPDATE`, `CREATE TABLE`, `DROP TABLE` akan menampilkan statistik ringkasan baris data yang terpengaruh (affected rows) atau pesan sukses.
4. **Log Error Terintegrasi**: Apabila terdapat kesalahan penulisan query, terminal akan mengembalikan pesan error orisinil dari PostgreSQL dalam blok notifikasi berwarna merah untuk mempermudah proses debugging.
5. **Penyelarasan GUI Otomatis**: Jika Anda menjalankan perintah SQL yang mengubah skema (seperti `CREATE TABLE`, `ALTER TABLE`) atau memodifikasi data baris, struktur GUI dashboard akan otomatis diperbarui.

---

## 🔒 8. Keamanan & Best Practices

1. **Gunakan Primary Key**: Selalu buat kolom `id` dengan tipe `SERIAL PRIMARY KEY` di setiap tabel. Ini sangat penting agar data dapat di-update atau dihapus secara akurat melalui antarmuka GUI.
2. **Validasi Tipe Data**: Pastikan input data baris sesuai dengan tipe data kolom (misalnya, jangan memasukkan teks ke dalam kolom bertipe `INTEGER`).
3. **Keamanan Kredensial**: Selalu gunakan kredensial spesifik pemilik database Anda saat melakukan integrasi di dalam aplikasi luar. Jangan pernah membocorkan kredensial `superadmin` PostgreSQL.

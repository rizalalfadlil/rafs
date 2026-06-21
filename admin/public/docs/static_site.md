# 🌐 Panduan Hosting Website Statis

Layanan Website Statis pada server **RAFS** memungkinkan Anda melakukan hosting berkas web berbasis HTML, CSS, JavaScript, serta aset gambar secara langsung. Halaman web yang di-hosting disajikan secara lokal dan terisolasi untuk mempermudah pengerjaan tugas kuliah maupun demonstrasi proyek secara mandiri.

---

## 🚀 1. Cara Kerja & URL Akses

- **Direktori Penyimpanan**: Semua berkas website Anda disimpan di folder `./www/<nama-situs>/` pada server.
- **Akses Publik**: Setiap website dapat diakses dari browser dengan alamat:
  `http://localhost:8080/sites/<nama-situs>/`
- **Syarat Keaktifan**: Agar website dapat diakses dengan benar, pastikan terdapat berkas **`index.html`** di direktori utama (root) folder website tersebut.

---

## 🏗️ 2. Metode Pembuatan Website

Ada dua metode mudah untuk menambahkan website statis baru ke dalam server RAFS:

### 2.1 Metode A: Clone dari Repositori Git

Metode ini sangat cocok jika kode sumber website Anda sudah disimpan di GitHub atau GitLab publik.

1. Masuk ke halaman **Static Sites**.
2. Klik tombol **Clone dari GitHub**.
3. Isi kolom input:
   - **URL Repositori**: Tautan repositori Git publik (contoh: `https://github.com/username/my-portfolio.git`).
   - **Nama Website**: Nama folder tujuan di server (contoh: `my-portfolio`).
4. Klik **Clone**. Server akan mengunduh kode proyek tersebut, menghapus folder metadata `.git` agar menghemat ruang penyimpanan, dan langsung menyajikannya.

### 2.2 Metode B: Upload File ZIP

Metode ini digunakan apabila file website berada di komputer lokal Anda.

1. Kompres seluruh file dan folder website Anda ke dalam satu file arsip **ZIP**.
   > [!IMPORTANT]
   > Pastikan file `index.html` berada di tingkat paling atas arsip ZIP (tidak terbungkus di dalam folder tambahan lagi di dalam ZIP).
2. Masuk ke halaman **Static Sites** dan klik **Unggah File ZIP**.
3. Masukkan **Nama Website** yang diinginkan.
4. Pilih file ZIP dari komputer Anda (ukuran file maksimal **50 MB**).
5. Klik **Unggah**. Server akan otomatis mengekstrak file ZIP tersebut ke direktori server.

---

## 🔍 3. Memantau Status Keaktifan (Active State)

Pada dashboard, setiap website statis akan ditandai dengan label status:

- **Active (Hijau)**: Berkas `index.html` berhasil ditemukan di root folder. Website siap dikunjungi dengan mengklik tombol **Buka Situs**.
- **Inactive (Abu-abu)**: File `index.html` tidak ditemukan di root folder (kemungkinan struktur folder di dalam ZIP salah atau repositori Git tidak memiliki index.html di root). Anda perlu merapikan struktur file di storage server agar website dapat diakses.

---

## 🗑️ 4. Menghapus Website

Jika suatu website sudah tidak digunakan:

1. Klik ikon **Hapus (Tong Sampah)** di samping nama website yang bersangkutan.
2. Konfirmasi penghapusan.
3. Server akan menghapus seluruh isi direktori website statis tersebut secara permanen.

---

## 🔒 5. Fitur Keamanan Ekstra

- **Validasi Nama Website**: Nama website hanya boleh menggunakan karakter alfanumerik, dash (`-`), dan underscore (`_`). Hal ini untuk mencegah serangan directory traversal pada URL.
- **Proteksi Zip Slip**: Server RAFS dilengkapi dengan validasi _Zip Slip_ saat proses dekompresi ZIP. Skrip pengekstrakan akan menolak dan membuang berkas jika terdeteksi mencoba menulis file ke luar batas direktori aman server (misalnya menggunakan nama file `../../target`).

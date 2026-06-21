# ☁️ Panduan Layanan Cloud Storage Pribadi

Layanan Cloud Storage pada server **RAFS** menyediakan penyimpanan awan privat untuk menyimpan, mengatur, mengunduh, dan membagikan berkas-berkas digital Anda secara praktis melalui antarmuka berbasis web.

---

## 🚀 1. Spesifikasi & Batasan Penyimpanan

- **Direktori Utama**: Berkas disimpan secara privat di dalam server pada folder `./storage_data/`.
- **Kuota Maksimal**: Kapasitas penyimpanan dibatasi maksimal sebesar **1 GB** (1024 MB) untuk memastikan penggunaan resource server lokal tetap efisien.
- **Ukuran Unggahan**: Batas maksimum ukuran berkas sekali unggah (upload batch) adalah **100 MB**.

---

## 📂 2. Fitur Utama Pengelola Berkas

### 2.1 Menavigasi & Membuat Folder

- **Eksplorasi Direktori**: Anda dapat mengklik nama folder pada daftar untuk masuk ke subdirektori, serta menggunakan tombol navigasi (breadcrumbs) untuk kembali ke folder di atasnya.
- **Membuat Folder Baru**:
  1. Klik tombol **Buat Folder**.
  2. Masukkan nama folder baru.
  3. Klik **Simpan**.
     > [!NOTE]
     > Nama item hanya boleh menggunakan huruf, angka, spasi, titik (`.`), hubung (`-`), dan garis bawah (`_`) untuk memastikan kompatibilitas sistem berkas lintas platform.

### 2.2 Mengunggah File (Multi-Upload)

Anda dapat mengunggah beberapa berkas sekaligus ke dalam direktori yang sedang dibuka:

1. Klik tombol **Unggah File**.
2. Pilih satu atau beberapa berkas dari komputer Anda.
3. Klik **Mulai Unggah**. Indikator progres akan menampilkan status pengunggahan hingga selesai.

### 2.3 Mengunduh & Pratinjau Berkas

- Untuk mengunduh berkas, klik tombol **Unduh** pada berkas yang bersangkutan.
- Server akan menyajikan berkas tersebut menggunakan mekanisme streaming native. Browser akan otomatis menampilkan pratinjau (preview) jika tipe berkas didukung (misalnya gambar, video, file teks, atau PDF), atau memicu dialog download jika tipe berkas biner.

### 2.4 Menghapus Berkas / Folder

- Anda dapat memilih satu atau beberapa file/folder dengan memberikan centang pada checkbox di daftar berkas.
- Klik tombol **Hapus Terpilih** untuk menghapusnya secara rekursif dari server secara permanen.
- _Demi keamanan, sistem memblokir tindakan penghapusan direktori root cloud storage._

---

## 🌐 3. Publikasi Berkas (Set Public)

Secara default, seluruh berkas dalam Cloud Storage bersifat **privat** dan hanya bisa diakses setelah pengguna masuk ke panel admin. Namun, Anda dapat mempublikasikan berkas agar dapat diakses oleh siapa saja tanpa autentikasi:

1. Pilih berkas atau folder privat yang ingin dibagikan.
2. Klik tombol **Publikasikan (Set Public)**.
3. Server akan menyalin berkas tersebut ke folder publik (`./public/`) secara rekursif.
4. Sistem akan mengembalikan tautan URL publik seperti:
   `/public/path-berkas/nama-berkas.jpg`
5. Tautan di atas dapat langsung diakses melalui URL:
   `http://localhost:8080/public/path-berkas/nama-berkas.jpg`

---

## 🔒 4. Keamanan & Pencegahan Kerentanan

- **Pencegahan Directory Traversal**: Server memvalidasi setiap path request menggunakan fungsi keamanan ketat (`getSafePath`). Fungsi ini mengevaluasi path absolut dan memblokir akses secara instan jika mendeteksi upaya manipulasi path (seperti penggunaan parameter `../` atau `..\\`) untuk mengintip berkas sistem di luar direktori `./storage_data/` atau `./public/`.
- **Sanitasi Nama Berkas**: Karakter nama file yang diunggah akan otomatis dibersihkan apabila mengandung simbol-simbol ilegal yang dapat mengganggu kestabilan sistem operasi server.

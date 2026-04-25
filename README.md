# School Nutrition Vault API

API backend untuk mengelola inventaris bahan/barang program MBG di sekolah-sekolah (stok, batas minimum, dan mutasi barang masuk/keluar).

## Tentang Proyek

Proyek ini dibuat untuk:
- Menyimpan data master barang (nama, kategori, stok, minimum stok, satuan).
- Mencatat transaksi stok masuk (`IN`) dan keluar (`OUT`).
- Mengupdate stok barang secara otomatis setiap transaksi.
- Menyediakan monitoring barang dengan stok kritis (`stock <= minimum_stock`).

## Teknologi yang Digunakan

- Go 1.25.5
- Gin (HTTP framework)
- PostgreSQL 15
- Docker & Docker Compose (opsional, sangat direkomendasikan)

## Struktur Singkat

- `main.go`: entry point aplikasi.
- `routes/routes.go`: definisi endpoint API.
- `controllers/`: logic CRUD item dan transaksi stok.
- `models/`: struktur data item dan transaction.
- `config/database.go`: koneksi PostgreSQL via environment variable.
- `ddl.sql`: inisialisasi tabel database.
- `School_Vault_API.postman_collection.json`: koleksi request untuk testing.

## Skema Database

Pada startup container PostgreSQL (via Docker Compose), file `ddl.sql` akan dijalankan otomatis dan membuat:
- Tabel `items`
- Tabel `transactions`

Relasi:
- `transactions.item_id` mereferensikan `items.id` (`ON DELETE CASCADE`).

## Cara Menjalankan Proyek

## Opsi 1 (Direkomendasikan): Jalankan dengan Docker Compose

### Prasyarat
- Docker
- Docker Compose

### Langkah
1. Pastikan berada di root project.
2. Jalankan:

```bash
docker compose up --build
```

3. Tunggu sampai log aplikasi menampilkan server berjalan di port `8080`.
4. API siap diakses di:

```text
http://localhost:8080
```

### Menghentikan

```bash
docker compose down
```

Jika ingin sekaligus menghapus volume database (reset data):

```bash
docker compose down -v
```

## Opsi 2: Jalankan Secara Lokal (Tanpa Docker)

### Prasyarat
- Go 1.25.5+ 
- PostgreSQL aktif di lokal

### 1. Buat database
- Nama default yang dipakai aplikasi: `mbg_inventory`

### 2. Jalankan script SQL
Import `ddl.sql` ke database `mbg_inventory`.

Contoh (opsional, sesuaikan user/password Anda):

```bash
psql -U postgres -d mbg_inventory -f ddl.sql
```

### 3. Set environment variable (opsional)
Jika tidak diset, aplikasi memakai default berikut:
- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=password123`
- `DB_NAME=mbg_inventory`

### 4. Jalankan aplikasi

```bash
go mod tidy
go run main.go
```

Server akan aktif di `http://localhost:8080`.

## Endpoint API

Base URL: `http://localhost:8080`

| Method | Endpoint | Fungsi |
|---|---|---|
| GET | `/ping` | Health check API |
| POST | `/items` | Tambah barang |
| GET | `/items` | Ambil semua barang |
| PUT | `/items/:id` | Ubah profil barang |
| DELETE | `/items/:id` | Hapus barang |
| GET | `/items/critical` | Ambil barang stok kritis |
| POST | `/transactions` | Catat transaksi IN/OUT |

## Format Request Penting

### 1) Tambah Barang
`POST /items`

```json
{
  "name": "Susu UHT Coklat 200ml",
  "category": "Susu",
  "stock": 50,
  "minimum_stock": 20,
  "unit": "Kotak"
}
```

### 2) Transaksi Stok Keluar
`POST /transactions`

```json
{
  "item_id": 1,
  "type": "OUT",
  "quantity": 40,
  "notes": "Diberikan ke siswa kelas 1 & 2"
}
```

### 3) Transaksi Stok Masuk
`POST /transactions`

```json
{
  "item_id": 1,
  "type": "IN",
  "quantity": 100,
  "notes": "Restok dari Dinas Kesehatan"
}
```

Catatan validasi transaksi:
- `type` wajib `IN` atau `OUT`.
- `OUT` akan gagal jika stok tidak mencukupi.

## Testing Menggunakan Postman

Project ini sudah menyediakan collection:
- `School_Vault_API.postman_collection.json`

### A. Import Collection
1. Buka Postman.
2. Klik **Import**.
3. Pilih file `School_Vault_API.postman_collection.json`.
4. Pastikan base URL yang dipakai adalah `http://localhost:8080`.

### B. Urutan Testing yang Disarankan
Jalankan request secara berurutan agar hasil valid:

1. `1. Ping Server`
- Expected: status `200`, message welcome.

2. `2. Tambah Barang Baru (Create Item)`
- Expected: status `201`, response berisi data barang + `id`.

3. `3. Lihat Semua Barang (Get Items)`
- Expected: status `200`, list barang muncul.

4. `4. Ubah Data Barang (Update Item)`
- Expected: status `200`, pesan update berhasil.

5. `5. Transaksi Barang Keluar (OUT)`
- Expected: status `201` jika stok cukup.
- Jika stok kurang: status `400` dengan pesan gagal.

6. `6. Transaksi Barang Masuk (IN)`
- Expected: status `201`, stok bertambah.

7. `7. Cek Stok Kritis (Critical Stocks)`
- Expected: status `200`, menampilkan item dengan stok <= minimum.

8. `8. Hapus Barang (Delete Item)`
- Expected: status `200` jika item ada.

### C. Tips Saat Testing
- Jika endpoint gagal konek, cek apakah server sudah aktif di port `8080`.
- Jika error database, cek service PostgreSQL dan kredensial environment variable.
- Untuk reset data saat Docker digunakan, jalankan `docker compose down -v` lalu `docker compose up --build`.

## Troubleshooting Singkat

- `connection refused` ke DB:
  - Pastikan PostgreSQL berjalan dan `DB_HOST/DB_PORT` benar.

- Gagal `OUT` transaksi:
  - Pastikan `quantity` tidak melebihi stok item.

- Tabel belum ada:
  - Pastikan `ddl.sql` sudah dieksekusi.

## Pengembangan Lanjutan (Opsional)

- Tambah endpoint riwayat transaksi (`GET /transactions`).
- Tambah pagination/filter untuk list item.
- Tambah autentikasi (JWT).
- Tambah automated test (unit/integration).

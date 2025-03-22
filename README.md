# API Specification Jubel - Aplikasi Jual Beli Barang Bekas Mahasiswa

## Informasi Umum

Base URL: `https://besbd.fuadfakhruz.id/api/v1`

### Autentikasi

Semua endpoint (kecuali yang disebutkan secara eksplisit) memerlukan autentikasi dengan JWT token. Token harus disertakan dalam header HTTP dengan format:

```
Authorization: Bearer <token>
```

### Format Response

#### Response Sukses

```json
{
  "status": "success",
  "message": "Pesan sukses",
  "data": {
    /* data response */
  }
}
```

#### Response Error

```json
{
  "status": "error",
  "message": "Pesan error",
  "errors": [
    /* detail error jika ada */
  ]
}
```

#### Response dengan Paginasi

```json
{
  "status": "success",
  "message": "Pesan sukses",
  "data": [
    /* array data */
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 100,
    "total_pages": 10
  }
}
```

## Endpoints

### Auth

#### Register

**Deskripsi**: Mendaftarkan pengguna baru.

- **URL**: `/auth/register`
- **Method**: `POST`
- **Auth Required**: Tidak
- **Body**:

```json
{
  "nama": "Budi Santoso",
  "email": "budi@example.com",
  "password": "password123",
  "no_hp": "081234567890",
  "alamat": "Jl. Sudirman No. 123, Jakarta"
}
```

- **Response Success (201)**:

```json
{
  "status": "success",
  "message": "Pendaftaran berhasil",
  "data": {
    "user_id": 1
  }
}
```

#### Login

**Deskripsi**: Melakukan login pengguna.

- **URL**: `/auth/login`
- **Method**: `POST`
- **Auth Required**: Tidak
- **Body**:

```json
{
  "email": "budi@example.com",
  "password": "password123"
}
```

- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Login berhasil",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "nama": "Budi Santoso",
      "email": "budi@example.com",
      "no_hp": "081234567890",
      "alamat": "Jl. Sudirman No. 123, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T10:00:00Z",
      "updated_at": "2025-03-23T10:00:00Z"
    }
  }
}
```

### Users

#### Get Current User

**Deskripsi**: Mendapatkan data pengguna yang sedang login.

- **URL**: `/users/me`
- **Method**: `GET`
- **Auth Required**: Ya
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Data pengguna berhasil diambil",
  "data": {
    "id": 1,
    "nama": "Budi Santoso",
    "email": "budi@example.com",
    "no_hp": "081234567890",
    "alamat": "Jl. Sudirman No. 123, Jakarta",
    "role": "user",
    "created_at": "2025-03-23T10:00:00Z",
    "updated_at": "2025-03-23T10:00:00Z"
  }
}
```

#### Get User by ID

**Deskripsi**: Mendapatkan data pengguna berdasarkan ID.

- **URL**: `/users/:id`
- **Method**: `GET`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID pengguna
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Data pengguna berhasil diambil",
  "data": {
    "id": 1,
    "nama": "Budi Santoso",
    "email": "budi@example.com",
    "no_hp": "081234567890",
    "alamat": "Jl. Sudirman No. 123, Jakarta",
    "role": "user",
    "created_at": "2025-03-23T10:00:00Z",
    "updated_at": "2025-03-23T10:00:00Z"
  }
}
```

#### Get All Users (Admin)

**Deskripsi**: Mendapatkan daftar semua pengguna (hanya admin).

- **URL**: `/users`
- **Method**: `GET`
- **Auth Required**: Ya (Role: Admin)
- **Query Params**:
  - `page` - Halaman (default: 1)
  - `limit` - Jumlah item per halaman (default: 10)
  - `search` - Kata kunci pencarian (optional)
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Daftar pengguna berhasil diambil",
  "data": [
    {
      "id": 1,
      "nama": "Budi Santoso",
      "email": "budi@example.com",
      "no_hp": "081234567890",
      "alamat": "Jl. Sudirman No. 123, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T10:00:00Z",
      "updated_at": "2025-03-23T10:00:00Z"
    },
    {
      "id": 2,
      "nama": "Andi Wijaya",
      "email": "andi@example.com",
      "no_hp": "087654321098",
      "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T11:00:00Z",
      "updated_at": "2025-03-23T11:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 2,
    "total_pages": 1
  }
}
```

#### Update User

**Deskripsi**: Memperbarui data pengguna.

- **URL**: `/users/:id`
- **Method**: `PATCH`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID pengguna
- **Body**:

```json
{
  "nama": "Budi Santoso Updated",
  "no_hp": "08123456789",
  "alamat": "Jl. Contoh Baru No. 456, Jakarta"
}
```

- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Data pengguna berhasil diperbarui",
  "data": {
    "id": 1,
    "nama": "Budi Santoso Updated",
    "email": "budi@example.com",
    "no_hp": "08123456789",
    "alamat": "Jl. Contoh Baru No. 456, Jakarta",
    "role": "user",
    "created_at": "2025-03-23T10:00:00Z",
    "updated_at": "2025-03-23T12:00:00Z"
  }
}
```

#### Delete User

**Deskripsi**: Menghapus pengguna (soft delete).

- **URL**: `/users/:id`
- **Method**: `DELETE`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID pengguna
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Pengguna berhasil dihapus",
  "data": null
}
```

### Items

#### Create Item

**Deskripsi**: Menambahkan barang baru.

- **URL**: `/items`
- **Method**: `POST`
- **Auth Required**: Ya
- **Body**:

```json
{
  "nama_barang": "Laptop Bekas",
  "harga": 3500000,
  "kategori": "Elektronik",
  "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB."
}
```

- **Response Success (201)**:

```json
{
  "status": "success",
  "message": "Barang berhasil ditambahkan",
  "data": {
    "id": 1,
    "penjual_id": 1,
    "nama_barang": "Laptop Bekas",
    "harga": 3500000,
    "kategori": "Elektronik",
    "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
    "gambar": null,
    "status": "Tersedia",
    "created_at": "2025-03-23T13:00:00Z",
    "penjual": {
      "id": 1,
      "nama": "Budi Santoso",
      "email": "budi@example.com",
      "no_hp": "081234567890",
      "alamat": "Jl. Sudirman No. 123, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T10:00:00Z",
      "updated_at": "2025-03-23T10:00:00Z"
    }
  }
}
```

#### Get Item by ID

**Deskripsi**: Mendapatkan data barang berdasarkan ID.

- **URL**: `/items/:id`
- **Method**: `GET`
- **Auth Required**: Tidak
- **URL Params**:
  - `id` - ID barang
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Data barang berhasil diambil",
  "data": {
    "id": 1,
    "penjual_id": 1,
    "nama_barang": "Laptop Bekas",
    "harga": 3500000,
    "kategori": "Elektronik",
    "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
    "gambar": "1_1711194000.jpg",
    "status": "Tersedia",
    "created_at": "2025-03-23T13:00:00Z",
    "penjual": {
      "id": 1,
      "nama": "Budi Santoso",
      "email": "budi@example.com",
      "no_hp": "081234567890",
      "alamat": "Jl. Sudirman No. 123, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T10:00:00Z",
      "updated_at": "2025-03-23T10:00:00Z"
    }
  }
}
```

#### Get All Items

**Deskripsi**: Mendapatkan daftar semua barang dengan filter.

- **URL**: `/items`
- **Method**: `GET`
- **Auth Required**: Tidak
- **Query Params**:
  - `page` - Halaman (default: 1)
  - `limit` - Jumlah item per halaman (default: 10)
  - `search` - Kata kunci pencarian (optional)
  - `kategori` - Filter berdasarkan kategori: "Buku", "Elektronik", "Perabotan", "Kos-kosan", "Lainnya" (optional)
  - `status` - Filter berdasarkan status: "Tersedia", "Terjual", "Dihapus" (optional)
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Daftar barang berhasil diambil",
  "data": [
    {
      "id": 1,
      "penjual_id": 1,
      "nama_barang": "Laptop Bekas",
      "harga": 3500000,
      "kategori": "Elektronik",
      "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
      "gambar": "1_1711194000.jpg",
      "status": "Tersedia",
      "created_at": "2025-03-23T13:00:00Z",
      "penjual": {
        "id": 1,
        "nama": "Budi Santoso",
        "email": "budi@example.com",
        "no_hp": "081234567890",
        "alamat": "Jl. Sudirman No. 123, Jakarta",
        "role": "user",
        "created_at": "2025-03-23T10:00:00Z",
        "updated_at": "2025-03-23T10:00:00Z"
      }
    },
    {
      "id": 2,
      "penjual_id": 2,
      "nama_barang": "Buku Algoritma",
      "harga": 85000,
      "kategori": "Buku",
      "deskripsi": "Buku Algoritma dan Pemrograman edisi terbaru.",
      "gambar": "2_1711197600.jpg",
      "status": "Tersedia",
      "created_at": "2025-03-23T14:00:00Z",
      "penjual": {
        "id": 2,
        "nama": "Andi Wijaya",
        "email": "andi@example.com",
        "no_hp": "087654321098",
        "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
        "role": "user",
        "created_at": "2025-03-23T11:00:00Z",
        "updated_at": "2025-03-23T11:00:00Z"
      }
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 2,
    "total_pages": 1
  }
}
```

#### Get Items by Seller

**Deskripsi**: Mendapatkan daftar barang berdasarkan penjual.

- **URL**: `/items/penjual/:id`
- **Method**: `GET`
- **Auth Required**: Tidak
- **URL Params**:
  - `id` - ID penjual
- **Query Params**:
  - `page` - Halaman (default: 1)
  - `limit` - Jumlah item per halaman (default: 10)
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Daftar barang berhasil diambil",
  "data": [
    {
      "id": 1,
      "penjual_id": 1,
      "nama_barang": "Laptop Bekas",
      "harga": 3500000,
      "kategori": "Elektronik",
      "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
      "gambar": "1_1711194000.jpg",
      "status": "Tersedia",
      "created_at": "2025-03-23T13:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### Get My Items

**Deskripsi**: Mendapatkan daftar barang milik pengguna yang login.

- **URL**: `/items/my`
- **Method**: `GET`
- **Auth Required**: Ya
- **Query Params**:
  - `page` - Halaman (default: 1)
  - `limit` - Jumlah item per halaman (default: 10)
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Daftar barang berhasil diambil",
  "data": [
    {
      "id": 1,
      "penjual_id": 1,
      "nama_barang": "Laptop Bekas",
      "harga": 3500000,
      "kategori": "Elektronik",
      "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
      "gambar": "1_1711194000.jpg",
      "status": "Tersedia",
      "created_at": "2025-03-23T13:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### Update Item

**Deskripsi**: Memperbarui data barang.

- **URL**: `/items/:id`
- **Method**: `PATCH`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID barang
- **Body**:

```json
{
  "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
  "harga": 3800000,
  "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet."
}
```

- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Data barang berhasil diperbarui",
  "data": {
    "id": 1,
    "penjual_id": 1,
    "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
    "harga": 3800000,
    "kategori": "Elektronik",
    "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
    "gambar": "1_1711194000.jpg",
    "status": "Tersedia",
    "created_at": "2025-03-23T13:00:00Z",
    "penjual": {
      "id": 1,
      "nama": "Budi Santoso",
      "email": "budi@example.com",
      "no_hp": "081234567890",
      "alamat": "Jl. Sudirman No. 123, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T10:00:00Z",
      "updated_at": "2025-03-23T10:00:00Z"
    }
  }
}
```

#### Update Item Status

**Deskripsi**: Memperbarui status barang.

- **URL**: `/items/:id/status`
- **Method**: `PATCH`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID barang
- **Body**:

```json
{
  "status": "Tersedia"
}
```

- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Status barang berhasil diperbarui",
  "data": null
}
```

#### Upload Item Image

**Deskripsi**: Mengupload gambar barang.

- **URL**: `/items/:id/upload`
- **Method**: `POST`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID barang
- **Body**:
  - Form-data dengan key `gambar` dan value berupa file gambar
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Gambar berhasil diupload",
  "data": {
    "file_name": "1_1711194000.jpg"
  }
}
```

#### Delete Item

**Deskripsi**: Menghapus barang (soft delete).

- **URL**: `/items/:id`
- **Method**: `DELETE`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID barang
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Barang berhasil dihapus",
  "data": null
}
```

### Transactions

#### Create Transaction

**Deskripsi**: Membuat transaksi baru.

- **URL**: `/transactions`
- **Method**: `POST`
- **Auth Required**: Ya
- **Body**:

```json
{
  "barang_id": 1
}
```

- **Response Success (201)**:

```json
{
  "status": "success",
  "message": "Transaksi berhasil dibuat",
  "data": {
    "id": 1,
    "barang_id": 1,
    "pembeli_id": 2,
    "tanggal_transaksi": "2025-03-23T15:00:00Z",
    "status_transaksi": "Pending",
    "created_at": "2025-03-23T15:00:00Z",
    "barang": {
      "id": 1,
      "penjual_id": 1,
      "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
      "harga": 3800000,
      "kategori": "Elektronik",
      "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
      "gambar": "1_1711194000.jpg",
      "status": "Terjual",
      "created_at": "2025-03-23T13:00:00Z"
    },
    "pembeli": {
      "id": 2,
      "nama": "Andi Wijaya",
      "email": "andi@example.com",
      "no_hp": "087654321098",
      "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T11:00:00Z",
      "updated_at": "2025-03-23T11:00:00Z"
    }
  }
}
```

#### Get Transaction by ID

**Deskripsi**: Mendapatkan data transaksi berdasarkan ID.

- **URL**: `/transactions/:id`
- **Method**: `GET`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID transaksi
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Data transaksi berhasil diambil",
  "data": {
    "id": 1,
    "barang_id": 1,
    "pembeli_id": 2,
    "tanggal_transaksi": "2025-03-23T15:00:00Z",
    "status_transaksi": "Pending",
    "created_at": "2025-03-23T15:00:00Z",
    "barang": {
      "id": 1,
      "penjual_id": 1,
      "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
      "harga": 3800000,
      "kategori": "Elektronik",
      "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
      "gambar": "1_1711194000.jpg",
      "status": "Terjual",
      "created_at": "2025-03-23T13:00:00Z"
    },
    "pembeli": {
      "id": 2,
      "nama": "Andi Wijaya",
      "email": "andi@example.com",
      "no_hp": "087654321098",
      "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T11:00:00Z",
      "updated_at": "2025-03-23T11:00:00Z"
    }
  }
}
```

#### Get All Transactions (Admin)

**Deskripsi**: Mendapatkan daftar semua transaksi (hanya admin).

- **URL**: `/transactions`
- **Method**: `GET`
- **Auth Required**: Ya (Role: Admin)
- **Query Params**:
  - `page` - Halaman (default: 1)
  - `limit` - Jumlah item per halaman (default: 10)
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Daftar transaksi berhasil diambil",
  "data": [
    {
      "id": 1,
      "barang_id": 1,
      "pembeli_id": 2,
      "tanggal_transaksi": "2025-03-23T15:00:00Z",
      "status_transaksi": "Pending",
      "created_at": "2025-03-23T15:00:00Z",
      "barang": {
        "id": 1,
        "penjual_id": 1,
        "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
        "harga": 3800000,
        "kategori": "Elektronik",
        "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
        "gambar": "1_1711194000.jpg",
        "status": "Terjual",
        "created_at": "2025-03-23T13:00:00Z"
      },
      "pembeli": {
        "id": 2,
        "nama": "Andi Wijaya",
        "email": "andi@example.com",
        "no_hp": "087654321098",
        "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
        "role": "user",
        "created_at": "2025-03-23T11:00:00Z",
        "updated_at": "2025-03-23T11:00:00Z"
      }
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### Get My Transactions as Buyer

**Deskripsi**: Mendapatkan daftar transaksi sebagai pembeli.

- **URL**: `/transactions/as-pembeli`
- **Method**: `GET`
- **Auth Required**: Ya
- **Query Params**:
  - `page` - Halaman (default: 1)
  - `limit` - Jumlah item per halaman (default: 10)
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Daftar transaksi berhasil diambil",
  "data": [
    {
      "id": 1,
      "barang_id": 1,
      "pembeli_id": 2,
      "tanggal_transaksi": "2025-03-23T15:00:00Z",
      "status_transaksi": "Pending",
      "created_at": "2025-03-23T15:00:00Z",
      "barang": {
        "id": 1,
        "penjual_id": 1,
        "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
        "harga": 3800000,
        "kategori": "Elektronik",
        "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
        "gambar": "1_1711194000.jpg",
        "status": "Terjual",
        "created_at": "2025-03-23T13:00:00Z"
      }
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### Get My Transactions as Seller

**Deskripsi**: Mendapatkan daftar transaksi sebagai penjual.

- **URL**: `/transactions/as-penjual`
- **Method**: `GET`
- **Auth Required**: Ya
- **Query Params**:
  - `page` - Halaman (default: 1)
  - `limit` - Jumlah item per halaman (default: 10)
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Daftar transaksi berhasil diambil",
  "data": [
    {
      "id": 1,
      "barang_id": 1,
      "pembeli_id": 2,
      "tanggal_transaksi": "2025-03-23T15:00:00Z",
      "status_transaksi": "Pending",
      "created_at": "2025-03-23T15:00:00Z",
      "barang": {
        "id": 1,
        "penjual_id": 1,
        "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
        "harga": 3800000,
        "kategori": "Elektronik",
        "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
        "gambar": "1_1711194000.jpg",
        "status": "Terjual",
        "created_at": "2025-03-23T13:00:00Z"
      },
      "pembeli": {
        "id": 2,
        "nama": "Andi Wijaya",
        "email": "andi@example.com",
        "no_hp": "087654321098",
        "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
        "role": "user",
        "created_at": "2025-03-23T11:00:00Z",
        "updated_at": "2025-03-23T11:00:00Z"
      }
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### Update Transaction Status

**Deskripsi**: Memperbarui status transaksi.

- **URL**: `/transactions/:id/status`
- **Method**: `PATCH`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID transaksi
- **Body**:

```json
{
  "status": "Selesai"
}
```

- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Status transaksi berhasil diperbarui",
  "data": null
}
```

#### Delete Transaction

**Deskripsi**: Menghapus transaksi.

- **URL**: `/transactions/:id`
- **Method**: `DELETE`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID transaksi
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Transaksi berhasil dihapus",
  "data": null
}
```

### Chats

#### Send Message

**Deskripsi**: Mengirim pesan chat baru.

- **URL**: `/chats`
- **Method**: `POST`
- **Auth Required**: Ya
- **Body**:

```json
{
  "penerima_id": 1,
  "barang_id": 1,
  "pesan": "Halo, apakah barang ini masih tersedia?"
}
```

- **Response Success (201)**:

```json
{
  "status": "success",
  "message": "Pesan berhasil dikirim",
  "data": {
    "id": 1,
    "pengirim_id": 2,
    "penerima_id": 1,
    "barang_id": 1,
    "pesan": "Halo, apakah barang ini masih tersedia?",
    "timestamp": "2025-03-23T16:00:00Z",
    "dibaca": false,
    "pengirim": {
      "id": 2,
      "nama": "Andi Wijaya",
      "email": "andi@example.com",
      "no_hp": "087654321098",
      "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T11:00:00Z",
      "updated_at": "2025-03-23T11:00:00Z"
    },
    "penerima": {
      "id": 1,
      "nama": "Budi Santoso",
      "email": "budi@example.com",
      "no_hp": "081234567890",
      "alamat": "Jl. Sudirman No. 123, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T10:00:00Z",
      "updated_at": "2025-03-23T10:00:00Z"
    },
    "barang": {
      "id": 1,
      "penjual_id": 1,
      "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
      "harga": 3800000,
      "kategori": "Elektronik",
      "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
      "gambar": "1_1711194000.jpg",
      "status": "Terjual",
      "created_at": "2025-03-23T13:00:00Z"
    }
  }
}
```

#### Get Chat by ID

**Deskripsi**: Mendapatkan data chat berdasarkan ID.

- **URL**: `/chats/:id`
- **Method**: `GET`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID chat
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Data chat berhasil diambil",
  "data": {
    "id": 1,
    "pengirim_id": 2,
    "penerima_id": 1,
    "barang_id": 1,
    "pesan": "Halo, apakah barang ini masih tersedia?",
    "timestamp": "2025-03-23T16:00:00Z",
    "dibaca": false,
    "pengirim": {
      "id": 2,
      "nama": "Andi Wijaya",
      "email": "andi@example.com",
      "no_hp": "087654321098",
      "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T11:00:00Z",
      "updated_at": "2025-03-23T11:00:00Z"
    },
    "penerima": {
      "id": 1,
      "nama": "Budi Santoso",
      "email": "budi@example.com",
      "no_hp": "081234567890",
      "alamat": "Jl. Sudirman No. 123, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T10:00:00Z",
      "updated_at": "2025-03-23T10:00:00Z"
    },
    "barang": {
      "id": 1,
      "penjual_id": 1,
      "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
      "harga": 3800000,
      "kategori": "Elektronik",
      "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
      "gambar": "1_1711194000.jpg",
      "status": "Terjual",
      "created_at": "2025-03-23T13:00:00Z"
    }
  }
}
```

#### Get Chats by Item

**Deskripsi**: Mendapatkan daftar chat berdasarkan ID barang.

- **URL**: `/chats/barang/:id`
- **Method**: `GET`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID barang
- **Query Params**:
  - `page` - Halaman (default: 1)
  - `limit` - Jumlah item per halaman (default: 10)
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Daftar chat berhasil diambil",
  "data": [
    {
      "id": 1,
      "pengirim_id": 2,
      "penerima_id": 1,
      "barang_id": 1,
      "pesan": "Halo, apakah barang ini masih tersedia?",
      "timestamp": "2025-03-23T16:00:00Z",
      "dibaca": false,
      "pengirim": {
        "id": 2,
        "nama": "Andi Wijaya",
        "email": "andi@example.com",
        "no_hp": "087654321098",
        "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
        "role": "user",
        "created_at": "2025-03-23T11:00:00Z",
        "updated_at": "2025-03-23T11:00:00Z"
      },
      "penerima": {
        "id": 1,
        "nama": "Budi Santoso",
        "email": "budi@example.com",
        "no_hp": "081234567890",
        "alamat": "Jl. Sudirman No. 123, Jakarta",
        "role": "user",
        "created_at": "2025-03-23T10:00:00Z",
        "updated_at": "2025-03-23T10:00:00Z"
      },
      "barang": {
        "id": 1,
        "penjual_id": 1,
        "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
        "harga": 3800000,
        "kategori": "Elektronik",
        "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
        "gambar": "1_1711194000.jpg",
        "status": "Terjual",
        "created_at": "2025-03-23T13:00:00Z"
      }
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 1,
    "total_pages": 1
  }
}
```

#### Get Conversation

**Deskripsi**: Mendapatkan percakapan antara dua pengguna terkait barang tertentu.

- **URL**: `/chats/conversation`
- **Method**: `GET`
- **Auth Required**: Ya
- **Query Params**:
  - `penerima_id` - ID penerima
  - `barang_id` - ID barang
  - `page` - Halaman (default: 1)
  - `limit` - Jumlah item per halaman (default: 20)
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Percakapan berhasil diambil",
  "data": [
    {
      "id": 1,
      "pengirim_id": 2,
      "penerima_id": 1,
      "barang_id": 1,
      "pesan": "Halo, apakah barang ini masih tersedia?",
      "timestamp": "2025-03-23T16:00:00Z",
      "dibaca": true,
      "pengirim": {
        "id": 2,
        "nama": "Andi Wijaya",
        "email": "andi@example.com",
        "no_hp": "087654321098",
        "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
        "role": "user",
        "created_at": "2025-03-23T11:00:00Z",
        "updated_at": "2025-03-23T11:00:00Z"
      },
      "penerima": {
        "id": 1,
        "nama": "Budi Santoso",
        "email": "budi@example.com",
        "no_hp": "081234567890",
        "alamat": "Jl. Sudirman No. 123, Jakarta",
        "role": "user",
        "created_at": "2025-03-23T10:00:00Z",
        "updated_at": "2025-03-23T10:00:00Z"
      },
      "barang": {
        "id": 1,
        "penjual_id": 1,
        "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
        "harga": 3800000,
        "kategori": "Elektronik",
        "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
        "gambar": "1_1711194000.jpg",
        "status": "Terjual",
        "created_at": "2025-03-23T13:00:00Z"
      }
    },
    {
      "id": 2,
      "pengirim_id": 1,
      "penerima_id": 2,
      "barang_id": 1,
      "pesan": "Iya, masih tersedia. Silakan ditawar.",
      "timestamp": "2025-03-23T16:15:00Z",
      "dibaca": false,
      "pengirim": {
        "id": 1,
        "nama": "Budi Santoso",
        "email": "budi@example.com",
        "no_hp": "081234567890",
        "alamat": "Jl. Sudirman No. 123, Jakarta",
        "role": "user",
        "created_at": "2025-03-23T10:00:00Z",
        "updated_at": "2025-03-23T10:00:00Z"
      },
      "penerima": {
        "id": 2,
        "nama": "Andi Wijaya",
        "email": "andi@example.com",
        "no_hp": "087654321098",
        "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
        "role": "user",
        "created_at": "2025-03-23T11:00:00Z",
        "updated_at": "2025-03-23T11:00:00Z"
      },
      "barang": {
        "id": 1,
        "penjual_id": 1,
        "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
        "harga": 3800000,
        "kategori": "Elektronik",
        "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet.",
        "gambar": "1_1711194000.jpg",
        "status": "Terjual",
        "created_at": "2025-03-23T13:00:00Z"
      }
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total_items": 2,
    "total_pages": 1
  }
}
```

#### Get Chat Partners

**Deskripsi**: Mendapatkan daftar partner chat.

- **URL**: `/chats/partners`
- **Method**: `GET`
- **Auth Required**: Ya
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Daftar partner chat berhasil diambil",
  "data": [
    {
      "id": 1,
      "nama": "Budi Santoso",
      "email": "budi@example.com",
      "no_hp": "081234567890",
      "alamat": "Jl. Sudirman No. 123, Jakarta",
      "role": "user",
      "created_at": "2025-03-23T10:00:00Z",
      "updated_at": "2025-03-23T10:00:00Z"
    }
  ]
}
```

#### Mark as Read

**Deskripsi**: Menandai chat sebagai telah dibaca.

- **URL**: `/chats/:id/read`
- **Method**: `PATCH`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID chat
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Pesan ditandai sebagai dibaca",
  "data": null
}
```

#### Delete Chat

**Deskripsi**: Menghapus chat.

- **URL**: `/chats/:id`
- **Method**: `DELETE`
- **Auth Required**: Ya
- **URL Params**:
  - `id` - ID chat
- **Response Success (200)**:

```json
{
  "status": "success",
  "message": "Chat berhasil dihapus",
  "data": null
}
```

## Status Codes

- `200 OK` - Permintaan berhasil
- `201 Created` - Data berhasil dibuat
- `400 Bad Request` - Permintaan tidak valid
- `401 Unauthorized` - Autentikasi gagal
- `403 Forbidden` - Akses ditolak
- `404 Not Found` - Data tidak ditemukan
- `500 Internal Server Error` - Terjadi kesalahan server

## Appendix

### Enum Values

#### Role

- `user` - Pengguna biasa
- `admin` - Administrator

#### Item Status

- `Tersedia` - Barang tersedia untuk dibeli
- `Terjual` - Barang sudah terjual
- `Dihapus` - Barang telah dihapus (soft delete)

#### Item Category

- `Buku` - Kategori buku
- `Elektronik` - Kategori elektronik
- `Perabotan` - Kategori perabotan
- `Kos-kosan` - Kategori kos-kosan
- `Lainnya` - Kategori lainnya

#### Transaction Status

- `Pending` - Transaksi sedang berlangsung
- `Selesai` - Transaksi telah selesai
- `Dibatalkan` - Transaksi dibatalkan

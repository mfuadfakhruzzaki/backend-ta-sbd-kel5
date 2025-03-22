# Dokumentasi API Jubel - Jual Beli Barang Bekas Mahasiswa

## Deskripsi Umum

API untuk aplikasi jual beli barang bekas khusus mahasiswa. Sistem ini memfasilitasi transaksi jual-beli antar mahasiswa secara langsung dalam satu platform yang aman dan mudah digunakan.

## URL Dasar

```
https://besbd.fuadfakhruz.id/api/v1
```

## Autentikasi

API menggunakan JWT (JSON Web Token) untuk autentikasi. Token diperoleh melalui endpoint login dan harus disertakan dalam header Authorization untuk mengakses endpoint yang memerlukan autentikasi.

Format header:

```
Authorization: Bearer {token}
```

## Format Respons

Semua respons API menggunakan format JSON dengan struktur berikut:

### Respons Sukses

```json
{
  "status": "success",
  "message": "Pesan sukses",
  "data": { ... }
}
```

### Respons Sukses dengan Paginasi

```json
{
  "status": "success",
  "message": "Pesan sukses",
  "data": [ ... ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_items": 100,
    "total_pages": 10
  }
}
```

### Respons Error

```json
{
  "status": "error",
  "message": "Pesan error",
  "errors": [ ... ]
}
```

## Endpoints

### Autentikasi

#### Register - Mendaftar pengguna baru

- **URL**: `/auth/register`
- **Method**: POST
- **Auth**: Tidak
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
- **Respons Sukses (201 Created)**:
  ```json
  {
    "status": "success",
    "message": "Pendaftaran berhasil",
    "data": {
      "user_id": 1
    }
  }
  ```

#### Login - Autentikasi pengguna

- **URL**: `/auth/login`
- **Method**: POST
- **Auth**: Tidak
- **Body**:
  ```json
  {
    "email": "budi@example.com",
    "password": "password123"
  }
  ```
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Status transaksi berhasil diperbarui",
    "data": null
  }
  ```

#### Delete Transaction - Menghapus transaksi

- **URL**: `/transactions/:id`
- **Method**: DELETE
- **Auth**: Ya (Token)
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Transaksi berhasil dihapus",
    "data": null
  }
  ```

### Chat

#### Send Message - Mengirim pesan chat baru

- **URL**: `/chats`
- **Method**: POST
- **Auth**: Ya (Token)
- **Body**:
  ```json
  {
    "penerima_id": 1,
    "barang_id": 1,
    "pesan": "Halo, apakah barang ini masih tersedia?"
  }
  ```
- **Respons Sukses (201 Created)**:
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
      "timestamp": "2023-04-01T12:00:00Z",
      "dibaca": false,
      "pengirim": {
        "id": 2,
        "nama": "Andi Wijaya",
        "email": "andi@example.com",
        "no_hp": "081234567891",
        "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
        "role": "user",
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      },
      "penerima": {
        "id": 1,
        "nama": "Budi Santoso",
        "email": "budi@example.com",
        "no_hp": "081234567890",
        "alamat": "Jl. Sudirman No. 123, Jakarta",
        "role": "user",
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      },
      "barang": {
        "id": 1,
        "penjual_id": 1,
        "nama_barang": "Laptop Bekas",
        "harga": 3500000,
        "kategori": "Elektronik",
        "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
        "gambar": "1_1680350400.jpg",
        "status": "Tersedia",
        "created_at": "2023-04-01T12:00:00Z"
      }
    }
  }
  ```

#### Get Chat by ID - Mendapatkan data chat berdasarkan ID

- **URL**: `/chats/:id`
- **Method**: GET
- **Auth**: Ya (Token)
- **Respons Sukses (200 OK)**:
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
      "timestamp": "2023-04-01T12:00:00Z",
      "dibaca": false,
      "pengirim": {
        "id": 2,
        "nama": "Andi Wijaya",
        "email": "andi@example.com",
        "no_hp": "081234567891",
        "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
        "role": "user",
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      },
      "penerima": {
        "id": 1,
        "nama": "Budi Santoso",
        "email": "budi@example.com",
        "no_hp": "081234567890",
        "alamat": "Jl. Sudirman No. 123, Jakarta",
        "role": "user",
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      },
      "barang": {
        "id": 1,
        "penjual_id": 1,
        "nama_barang": "Laptop Bekas",
        "harga": 3500000,
        "kategori": "Elektronik",
        "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
        "gambar": "1_1680350400.jpg",
        "status": "Tersedia",
        "created_at": "2023-04-01T12:00:00Z"
      }
    }
  }
  ```

#### Get Chats by Item - Mendapatkan daftar chat berdasarkan ID barang

- **URL**: `/chats/barang/:id?page=1&limit=10`
- **Method**: GET
- **Auth**: Ya (Token)
- **Parameter Query**:
  - `page`: Nomor halaman (default: 1)
  - `limit`: Jumlah item per halaman (default: 10)
- **Respons Sukses (200 OK)**:
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
        "timestamp": "2023-04-01T12:00:00Z",
        "dibaca": false,
        "pengirim": {
          "id": 2,
          "nama": "Andi Wijaya",
          "email": "andi@example.com",
          "no_hp": "081234567891",
          "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
          "role": "user",
          "created_at": "2023-04-01T12:00:00Z",
          "updated_at": "2023-04-01T12:00:00Z"
        },
        "penerima": {
          "id": 1,
          "nama": "Budi Santoso",
          "email": "budi@example.com",
          "no_hp": "081234567890",
          "alamat": "Jl. Sudirman No. 123, Jakarta",
          "role": "user",
          "created_at": "2023-04-01T12:00:00Z",
          "updated_at": "2023-04-01T12:00:00Z"
        },
        "barang": {
          "id": 1,
          "penjual_id": 1,
          "nama_barang": "Laptop Bekas",
          "harga": 3500000,
          "kategori": "Elektronik",
          "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
          "gambar": "1_1680350400.jpg",
          "status": "Tersedia",
          "created_at": "2023-04-01T12:00:00Z"
        }
      }
      // ... data chat lainnya
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_items": 20,
      "total_pages": 2
    }
  }
  ```

#### Get Conversation - Mendapatkan percakapan antara dua pengguna terkait barang tertentu

- **URL**: `/chats/conversation?penerima_id=1&barang_id=1&page=1&limit=20`
- **Method**: GET
- **Auth**: Ya (Token)
- **Parameter Query**:
  - `penerima_id`: ID penerima pesan
  - `barang_id`: ID barang
  - `page`: Nomor halaman (default: 1)
  - `limit`: Jumlah item per halaman (default: 20)
- **Respons Sukses (200 OK)**:
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
        "timestamp": "2023-04-01T12:00:00Z",
        "dibaca": true,
        "pengirim": {
          "id": 2,
          "nama": "Andi Wijaya",
          "email": "andi@example.com",
          "no_hp": "081234567891",
          "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
          "role": "user",
          "created_at": "2023-04-01T12:00:00Z",
          "updated_at": "2023-04-01T12:00:00Z"
        },
        "penerima": {
          "id": 1,
          "nama": "Budi Santoso",
          "email": "budi@example.com",
          "no_hp": "081234567890",
          "alamat": "Jl. Sudirman No. 123, Jakarta",
          "role": "user",
          "created_at": "2023-04-01T12:00:00Z",
          "updated_at": "2023-04-01T12:00:00Z"
        },
        "barang": {
          "id": 1,
          "penjual_id": 1,
          "nama_barang": "Laptop Bekas",
          "harga": 3500000,
          "kategori": "Elektronik",
          "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
          "gambar": "1_1680350400.jpg",
          "status": "Tersedia",
          "created_at": "2023-04-01T12:00:00Z"
        }
      },
      {
        "id": 2,
        "pengirim_id": 1,
        "penerima_id": 2,
        "barang_id": 1,
        "pesan": "Ya, masih tersedia. Apakah Anda tertarik?",
        "timestamp": "2023-04-01T12:05:00Z",
        "dibaca": true,
        "pengirim": {
          "id": 1,
          "nama": "Budi Santoso",
          "email": "budi@example.com",
          "no_hp": "081234567890",
          "alamat": "Jl. Sudirman No. 123, Jakarta",
          "role": "user",
          "created_at": "2023-04-01T12:00:00Z",
          "updated_at": "2023-04-01T12:00:00Z"
        },
        "penerima": {
          "id": 2,
          "nama": "Andi Wijaya",
          "email": "andi@example.com",
          "no_hp": "081234567891",
          "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
          "role": "user",
          "created_at": "2023-04-01T12:00:00Z",
          "updated_at": "2023-04-01T12:00:00Z"
        },
        "barang": {
          "id": 1,
          "penjual_id": 1,
          "nama_barang": "Laptop Bekas",
          "harga": 3500000,
          "kategori": "Elektronik",
          "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
          "gambar": "1_1680350400.jpg",
          "status": "Tersedia",
          "created_at": "2023-04-01T12:00:00Z"
        }
      }
      // ... data chat lainnya
    ],
    "meta": {
      "page": 1,
      "limit": 20,
      "total_items": 5,
      "total_pages": 1
    }
  }
  ```

#### Get Chat Partners - Mendapatkan daftar partner chat

- **URL**: `/chats/partners`
- **Method**: GET
- **Auth**: Ya (Token)
- **Respons Sukses (200 OK)**:
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
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      },
      {
        "id": 3,
        "nama": "Citra Dewi",
        "email": "citra@example.com",
        "no_hp": "081234567892",
        "alamat": "Jl. MH Thamrin No. 789, Jakarta",
        "role": "user",
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      }
      // ... data partner chat lainnya
    ]
  }
  ```

#### Mark as Read - Menandai chat sebagai telah dibaca

- **URL**: `/chats/:id/read`
- **Method**: PATCH
- **Auth**: Ya (Token)
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Pesan ditandai sebagai dibaca",
    "data": null
  }
  ```

#### Delete Chat - Menghapus chat

- **URL**: `/chats/:id`
- **Method**: DELETE
- **Auth**: Ya (Token)
- **Respons Sukses (200 OK)**:

  ```json
  {
    "status": "success",
    "message": "Chat berhasil dihapus",
    "data": null
  }
  "status": "success",
  "message": "Login berhasil",
  "data": {
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "nama": "Budi Santoso",
    "email": "budi@example.com",
    "no_hp": "081234567890",
    "alamat": "Jl. Sudirman No. 123 Jakarta",
    "role": "user",
    "created_at": "2023-04-01T12:00:00Z",
    "updated_at": "2023-04-01T12:00:00Z"
    }
  }
  ```

### Pengguna (Users)

#### Get Current User - Mendapatkan data pengguna yang sedang login

- **URL**: `/users/me`
- **Method**: GET
- **Auth**: Ya (Token)
- **Respons Sukses (200 OK)**:
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
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    }
  }
  ```

#### Get User by ID - Mendapatkan data pengguna berdasarkan ID

- **URL**: `/users/:id`
- **Method**: GET
- **Auth**: Ya (Token)
- **Respons Sukses (200 OK)**:
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
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    }
  }
  ```

#### Get All Users - Mendapatkan daftar semua pengguna (Admin)

- **URL**: `/users?page=1&limit=10&search=`
- **Method**: GET
- **Auth**: Ya (Token Admin)
- **Parameter Query**:
  - `page`: Nomor halaman (default: 1)
  - `limit`: Jumlah item per halaman (default: 10)
  - `search`: Kata kunci pencarian
- **Respons Sukses (200 OK)**:
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
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      }
      // ... data pengguna lainnya
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_items": 100,
      "total_pages": 10
    }
  }
  ```

#### Update User - Memperbarui data pengguna

- **URL**: `/users/:id`
- **Method**: PATCH
- **Auth**: Ya (Token)
- **Body**:
  ```json
  {
    "nama": "Budi Santoso Updated",
    "no_hp": "08123456789",
    "alamat": "Jl. Contoh Baru No. 456, Jakarta"
  }
  ```
- **Respons Sukses (200 OK)**:
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
      "created_at": "2023-04-01T12:00:00Z",
      "updated_at": "2023-04-01T12:00:00Z"
    }
  }
  ```

#### Delete User - Menghapus pengguna (soft delete)

- **URL**: `/users/:id`
- **Method**: DELETE
- **Auth**: Ya (Token)
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Pengguna berhasil dihapus",
    "data": null
  }
  ```

### Barang (Items)

#### Create Item - Menambahkan barang baru

- **URL**: `/items`
- **Method**: POST
- **Auth**: Ya (Token)
- **Body**:
  ```json
  {
    "nama_barang": "Laptop Bekas",
    "harga": 3500000,
    "kategori": "Elektronik",
    "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB."
  }
  ```
- **Respons Sukses (201 Created)**:
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
      "created_at": "2023-04-01T12:00:00Z",
      "penjual": {
        "id": 1,
        "nama": "Budi Santoso",
        "email": "budi@example.com",
        "no_hp": "081234567890",
        "alamat": "Jl. Sudirman No. 123, Jakarta",
        "role": "user",
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      }
    }
  }
  ```

#### Get Item by ID - Mendapatkan data barang berdasarkan ID

- **URL**: `/items/:id`
- **Method**: GET
- **Auth**: Tidak
- **Respons Sukses (200 OK)**:
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
      "gambar": "1_1680350400.jpg",
      "status": "Tersedia",
      "created_at": "2023-04-01T12:00:00Z",
      "penjual": {
        "id": 1,
        "nama": "Budi Santoso",
        "email": "budi@example.com",
        "no_hp": "081234567890",
        "alamat": "Jl. Sudirman No. 123, Jakarta",
        "role": "user",
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      }
    }
  }
  ```

#### Get All Items - Mendapatkan daftar semua barang dengan filter

- **URL**: `/items?page=1&limit=10&search=&kategori=&status=`
- **Method**: GET
- **Auth**: Tidak
- **Parameter Query**:
  - `page`: Nomor halaman (default: 1)
  - `limit`: Jumlah item per halaman (default: 10)
  - `search`: Kata kunci pencarian
  - `kategori`: Filter berdasarkan kategori (Buku, Elektronik, Perabotan, Kos-kosan, Lainnya)
  - `status`: Filter berdasarkan status (Tersedia, Terjual, Dihapus)
- **Respons Sukses (200 OK)**:
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
        "gambar": "1_1680350400.jpg",
        "status": "Tersedia",
        "created_at": "2023-04-01T12:00:00Z",
        "penjual": {
          "id": 1,
          "nama": "Budi Santoso",
          "email": "budi@example.com",
          "no_hp": "081234567890",
          "alamat": "Jl. Sudirman No. 123, Jakarta",
          "role": "user",
          "created_at": "2023-04-01T12:00:00Z",
          "updated_at": "2023-04-01T12:00:00Z"
        }
      }
      // ... data barang lainnya
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_items": 100,
      "total_pages": 10
    }
  }
  ```

#### Get Items by Seller - Mendapatkan daftar barang berdasarkan penjual

- **URL**: `/items/penjual/:id?page=1&limit=10`
- **Method**: GET
- **Auth**: Tidak
- **Parameter Query**:
  - `page`: Nomor halaman (default: 1)
  - `limit`: Jumlah item per halaman (default: 10)
- **Respons Sukses (200 OK)**:
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
        "gambar": "1_1680350400.jpg",
        "status": "Tersedia",
        "created_at": "2023-04-01T12:00:00Z"
      }
      // ... data barang lainnya
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_items": 5,
      "total_pages": 1
    }
  }
  ```

#### Get My Items - Mendapatkan daftar barang milik pengguna yang login

- **URL**: `/items/my?page=1&limit=10`
- **Method**: GET
- **Auth**: Ya (Token)
- **Parameter Query**:
  - `page`: Nomor halaman (default: 1)
  - `limit`: Jumlah item per halaman (default: 10)
- **Respons Sukses (200 OK)**:
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
        "gambar": "1_1680350400.jpg",
        "status": "Tersedia",
        "created_at": "2023-04-01T12:00:00Z"
      }
      // ... data barang lainnya
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_items": 5,
      "total_pages": 1
    }
  }
  ```

#### Update Item - Memperbarui data barang

- **URL**: `/items/:id`
- **Method**: PATCH
- **Auth**: Ya (Token)
- **Body**:
  ```json
  {
    "nama_barang": "Laptop Bekas HP EliteBook 840 G3",
    "harga": 3800000,
    "deskripsi": "Laptop bekas HP EliteBook 840 G3 dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB. Baterai masih awet."
  }
  ```
- **Respons Sukses (200 OK)**:
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
      "gambar": "1_1680350400.jpg",
      "status": "Tersedia",
      "created_at": "2023-04-01T12:00:00Z",
      "penjual": {
        "id": 1,
        "nama": "Budi Santoso",
        "email": "budi@example.com",
        "no_hp": "081234567890",
        "alamat": "Jl. Sudirman No. 123, Jakarta",
        "role": "user",
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      }
    }
  }
  ```

#### Update Item Status - Memperbarui status barang

- **URL**: `/items/:id/status`
- **Method**: PATCH
- **Auth**: Ya (Token)
- **Body**:
  ```json
  {
    "status": "Tersedia"
  }
  ```
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Status barang berhasil diperbarui",
    "data": null
  }
  ```

#### Upload Item Image - Mengupload gambar barang

- **URL**: `/items/:id/upload`
- **Method**: POST
- **Auth**: Ya (Token)
- **Body**: Form-data dengan key `gambar` dan value file gambar
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Gambar berhasil diupload",
    "data": {
      "file_name": "1_1680350400.jpg"
    }
  }
  ```

#### Delete Item - Menghapus barang (soft delete)

- **URL**: `/items/:id`
- **Method**: DELETE
- **Auth**: Ya (Token)
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Barang berhasil dihapus",
    "data": null
  }
  ```

### Transaksi (Transactions)

#### Create Transaction - Membuat transaksi baru

- **URL**: `/transactions`
- **Method**: POST
- **Auth**: Ya (Token)
- **Body**:
  ```json
  {
    "barang_id": 1
  }
  ```
- **Respons Sukses (201 Created)**:
  ```json
  {
    "status": "success",
    "message": "Transaksi berhasil dibuat",
    "data": {
      "id": 1,
      "barang_id": 1,
      "pembeli_id": 2,
      "tanggal_transaksi": "2023-04-01T12:00:00Z",
      "status_transaksi": "Pending",
      "created_at": "2023-04-01T12:00:00Z",
      "barang": {
        "id": 1,
        "penjual_id": 1,
        "nama_barang": "Laptop Bekas",
        "harga": 3500000,
        "kategori": "Elektronik",
        "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
        "gambar": "1_1680350400.jpg",
        "status": "Terjual",
        "created_at": "2023-04-01T12:00:00Z"
      },
      "pembeli": {
        "id": 2,
        "nama": "Andi Wijaya",
        "email": "andi@example.com",
        "no_hp": "081234567891",
        "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
        "role": "user",
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      }
    }
  }
  ```

#### Get Transaction by ID - Mendapatkan data transaksi berdasarkan ID

- **URL**: `/transactions/:id`
- **Method**: GET
- **Auth**: Ya (Token)
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Data transaksi berhasil diambil",
    "data": {
      "id": 1,
      "barang_id": 1,
      "pembeli_id": 2,
      "tanggal_transaksi": "2023-04-01T12:00:00Z",
      "status_transaksi": "Pending",
      "created_at": "2023-04-01T12:00:00Z",
      "barang": {
        "id": 1,
        "penjual_id": 1,
        "nama_barang": "Laptop Bekas",
        "harga": 3500000,
        "kategori": "Elektronik",
        "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
        "gambar": "1_1680350400.jpg",
        "status": "Terjual",
        "created_at": "2023-04-01T12:00:00Z"
      },
      "pembeli": {
        "id": 2,
        "nama": "Andi Wijaya",
        "email": "andi@example.com",
        "no_hp": "081234567891",
        "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
        "role": "user",
        "created_at": "2023-04-01T12:00:00Z",
        "updated_at": "2023-04-01T12:00:00Z"
      }
    }
  }
  ```

#### Get All Transactions - Mendapatkan daftar semua transaksi (Admin)

- **URL**: `/transactions?page=1&limit=10`
- **Method**: GET
- **Auth**: Ya (Token Admin)
- **Parameter Query**:
  - `page`: Nomor halaman (default: 1)
  - `limit`: Jumlah item per halaman (default: 10)
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Daftar transaksi berhasil diambil",
    "data": [
      {
        "id": 1,
        "barang_id": 1,
        "pembeli_id": 2,
        "tanggal_transaksi": "2023-04-01T12:00:00Z",
        "status_transaksi": "Pending",
        "created_at": "2023-04-01T12:00:00Z",
        "barang": {
          "id": 1,
          "penjual_id": 1,
          "nama_barang": "Laptop Bekas",
          "harga": 3500000,
          "kategori": "Elektronik",
          "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
          "gambar": "1_1680350400.jpg",
          "status": "Terjual",
          "created_at": "2023-04-01T12:00:00Z"
        },
        "pembeli": {
          "id": 2,
          "nama": "Andi Wijaya",
          "email": "andi@example.com",
          "no_hp": "081234567891",
          "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
          "role": "user",
          "created_at": "2023-04-01T12:00:00Z",
          "updated_at": "2023-04-01T12:00:00Z"
        }
      }
      // ... data transaksi lainnya
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_items": 50,
      "total_pages": 5
    }
  }
  ```

#### Get My Transactions as Buyer - Mendapatkan daftar transaksi sebagai pembeli

- **URL**: `/transactions/as-pembeli?page=1&limit=10`
- **Method**: GET
- **Auth**: Ya (Token)
- **Parameter Query**:
  - `page`: Nomor halaman (default: 1)
  - `limit`: Jumlah item per halaman (default: 10)
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Daftar transaksi berhasil diambil",
    "data": [
      {
        "id": 1,
        "barang_id": 1,
        "pembeli_id": 2,
        "tanggal_transaksi": "2023-04-01T12:00:00Z",
        "status_transaksi": "Pending",
        "created_at": "2023-04-01T12:00:00Z",
        "barang": {
          "id": 1,
          "penjual_id": 1,
          "nama_barang": "Laptop Bekas",
          "harga": 3500000,
          "kategori": "Elektronik",
          "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
          "gambar": "1_1680350400.jpg",
          "status": "Terjual",
          "created_at": "2023-04-01T12:00:00Z"
        }
      }
      // ... data transaksi lainnya
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_items": 5,
      "total_pages": 1
    }
  }
  ```

#### Get My Transactions as Seller - Mendapatkan daftar transaksi sebagai penjual

- **URL**: `/transactions/as-penjual?page=1&limit=10`
- **Method**: GET
- **Auth**: Ya (Token)
- **Parameter Query**:
  - `page`: Nomor halaman (default: 1)
  - `limit`: Jumlah item per halaman (default: 10)
- **Respons Sukses (200 OK)**:
  ```json
  {
    "status": "success",
    "message": "Daftar transaksi berhasil diambil",
    "data": [
      {
        "id": 1,
        "barang_id": 1,
        "pembeli_id": 2,
        "tanggal_transaksi": "2023-04-01T12:00:00Z",
        "status_transaksi": "Pending",
        "created_at": "2023-04-01T12:00:00Z",
        "barang": {
          "id": 1,
          "penjual_id": 1,
          "nama_barang": "Laptop Bekas",
          "harga": 3500000,
          "kategori": "Elektronik",
          "deskripsi": "Laptop bekas HP EliteBook dalam kondisi baik. Spesifikasi: Core i5, RAM 8GB, SSD 256GB.",
          "gambar": "1_1680350400.jpg",
          "status": "Terjual",
          "created_at": "2023-04-01T12:00:00Z"
        },
        "pembeli": {
          "id": 2,
          "nama": "Andi Wijaya",
          "email": "andi@example.com",
          "no_hp": "081234567891",
          "alamat": "Jl. Gatot Subroto No. 456, Jakarta",
          "role": "user",
          "created_at": "2023-04-01T12:00:00Z",
          "updated_at": "2023-04-01T12:00:00Z"
        }
      }
      // ... data transaksi lainnya
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_items": 8,
      "total_pages": 1
    }
  }
  ```

#### Update Transaction Status - Memperbarui status transaksi

- **URL**: `/transactions/:id/status`
- **Method**: PATCH
- **Auth**: Ya (Token)
- **Body**:
  ```json
  {
    "status": "Selesai"
  }
  ```
- **Respons Sukses (200 OK)**:
  ```json
  {
  ```

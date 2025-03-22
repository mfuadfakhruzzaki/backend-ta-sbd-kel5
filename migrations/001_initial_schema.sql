-- Buat enum type untuk kolom yang membutuhkannya
CREATE TYPE user_role AS ENUM ('user', 'admin');
CREATE TYPE item_status AS ENUM ('Tersedia', 'Terjual', 'Dihapus');
CREATE TYPE item_category AS ENUM ('Buku', 'Elektronik', 'Perabotan', 'Kos-kosan', 'Lainnya');
CREATE TYPE transaction_status AS ENUM ('Pending', 'Selesai', 'Dibatalkan');

-- Tabel pengguna
CREATE TABLE pengguna (
    id SERIAL PRIMARY KEY,
    nama VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    no_hp VARCHAR(15) NOT NULL,
    alamat TEXT,
    role user_role DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Buat index untuk pencarian email
CREATE INDEX idx_pengguna_email ON pengguna(email);

-- Tabel barang
CREATE TABLE barang (
    id SERIAL PRIMARY KEY,
    penjual_id INT NOT NULL,
    nama_barang VARCHAR(100) NOT NULL,
    harga DECIMAL(10,2) NOT NULL,
    kategori item_category NOT NULL,
    deskripsi TEXT,
    gambar VARCHAR(255),
    status item_status DEFAULT 'Tersedia',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    FOREIGN KEY (penjual_id) REFERENCES pengguna(id) ON DELETE CASCADE
);

-- Buat index untuk pencarian dan filter
CREATE INDEX idx_barang_penjual ON barang(penjual_id);
CREATE INDEX idx_barang_kategori ON barang(kategori);
CREATE INDEX idx_barang_status ON barang(status);

-- Tabel transaksi
CREATE TABLE transaksi (
    id SERIAL PRIMARY KEY,
    barang_id INT NOT NULL,
    pembeli_id INT NOT NULL,
    tanggal_transaksi TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status_transaksi transaction_status DEFAULT 'Pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (barang_id) REFERENCES barang(id) ON DELETE CASCADE,
    FOREIGN KEY (pembeli_id) REFERENCES pengguna(id) ON DELETE CASCADE
);

-- Buat index untuk pencarian transaksi
CREATE INDEX idx_transaksi_barang ON transaksi(barang_id);
CREATE INDEX idx_transaksi_pembeli ON transaksi(pembeli_id);

-- Tabel chat
CREATE TABLE chat (
    id SERIAL PRIMARY KEY,
    pengirim_id INT NOT NULL,
    penerima_id INT NOT NULL,
    barang_id INT NOT NULL,
    pesan TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    dibaca BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (pengirim_id) REFERENCES pengguna(id) ON DELETE CASCADE,
    FOREIGN KEY (penerima_id) REFERENCES pengguna(id) ON DELETE CASCADE,
    FOREIGN KEY (barang_id) REFERENCES barang(id) ON DELETE CASCADE
);

-- Buat index untuk pencarian chat
CREATE INDEX idx_chat_pengirim ON chat(pengirim_id);
CREATE INDEX idx_chat_penerima ON chat(penerima_id);
CREATE INDEX idx_chat_barang ON chat(barang_id);
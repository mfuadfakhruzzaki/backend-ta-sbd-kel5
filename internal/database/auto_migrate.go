package database

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// AutoMigrate melakukan migrasi database otomatis saat aplikasi pertama kali dijalankan
func AutoMigrate(db *gorm.DB) error {
	log.Println("Checking database schema...")

	// Cek apakah tabel pengguna sudah ada
	if !tableExists(db, "pengguna") {
		log.Println("Tables not found. Running auto-migration...")
		
		// Buat tipe enum
		if err := createEnumTypes(db); err != nil {
			return fmt.Errorf("failed to create ENUM types: %w", err)
		}
		
		// Buat tabel-tabel
		if err := createTables(db); err != nil {
			return fmt.Errorf("failed to create tables: %w", err)
		}
		
		// Buat indeks
		if err := createIndexes(db); err != nil {
			return fmt.Errorf("failed to create indexes: %w", err)
		}
		
		log.Println("Auto-migration completed successfully.")
	} else {
		log.Println("Database schema already exists. Skipping auto-migration.")
	}
	
	return nil
}

// tableExists memeriksa apakah tabel sudah ada di database
func tableExists(db *gorm.DB, tableName string) bool {
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = ?", tableName).Count(&count)
	return count > 0
}

// createEnumTypes membuat tipe enum yang dibutuhkan
func createEnumTypes(db *gorm.DB) error {
	// Buat tipe enum untuk user_role
	if err := db.Exec(`CREATE TYPE user_role AS ENUM ('user', 'admin');`).Error; err != nil {
		return err
	}

	// Buat tipe enum untuk item_status
	if err := db.Exec(`CREATE TYPE item_status AS ENUM ('Tersedia', 'Terjual', 'Dihapus');`).Error; err != nil {
		return err
	}

	// Buat tipe enum untuk item_category
	if err := db.Exec(`CREATE TYPE item_category AS ENUM ('Buku', 'Elektronik', 'Perabotan', 'Kos-kosan', 'Lainnya');`).Error; err != nil {
		return err
	}

	// Buat tipe enum untuk transaction_status
	if err := db.Exec(`CREATE TYPE transaction_status AS ENUM ('Pending', 'Selesai', 'Dibatalkan');`).Error; err != nil {
		return err
	}

	return nil
}

// createTables membuat tabel-tabel database
func createTables(db *gorm.DB) error {
	// Buat tabel pengguna
	if err := db.Exec(`
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
	`).Error; err != nil {
		return err
	}

	// Buat tabel barang
	if err := db.Exec(`
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
	`).Error; err != nil {
		return err
	}

	// Buat tabel transaksi
	if err := db.Exec(`
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
	`).Error; err != nil {
		return err
	}

	// Buat tabel chat
	if err := db.Exec(`
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
	`).Error; err != nil {
		return err
	}

	return nil
}

// createIndexes membuat indeks untuk meningkatkan performa query
func createIndexes(db *gorm.DB) error {
	// Indeks untuk pencarian email
	if err := db.Exec(`CREATE INDEX idx_pengguna_email ON pengguna(email);`).Error; err != nil {
		return err
	}

	// Indeks untuk tabel barang
	if err := db.Exec(`CREATE INDEX idx_barang_penjual ON barang(penjual_id);`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE INDEX idx_barang_kategori ON barang(kategori);`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE INDEX idx_barang_status ON barang(status);`).Error; err != nil {
		return err
	}

	// Indeks untuk tabel transaksi
	if err := db.Exec(`CREATE INDEX idx_transaksi_barang ON transaksi(barang_id);`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE INDEX idx_transaksi_pembeli ON transaksi(pembeli_id);`).Error; err != nil {
		return err
	}

	// Indeks untuk tabel chat
	if err := db.Exec(`CREATE INDEX idx_chat_pengirim ON chat(pengirim_id);`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE INDEX idx_chat_penerima ON chat(penerima_id);`).Error; err != nil {
		return err
	}
	if err := db.Exec(`CREATE INDEX idx_chat_barang ON chat(barang_id);`).Error; err != nil {
		return err
	}

	return nil
}
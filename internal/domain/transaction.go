package domain

import (
	"time"
)

// Status transaksi
type TransactionStatus string

const (
	StatusPending   TransactionStatus = "Pending"
	StatusSelesai   TransactionStatus = "Selesai"
	StatusDibatalkan TransactionStatus = "Dibatalkan"
)

// Transaction merepresentasikan transaksi jual beli
type Transaction struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	BarangID        uint              `gorm:"column:barang_id;not null" json:"barang_id"`
	PembeliID       uint              `gorm:"column:pembeli_id;not null" json:"pembeli_id"`
	TanggalTransaksi time.Time         `gorm:"column:tanggal_transaksi;default:CURRENT_TIMESTAMP" json:"tanggal_transaksi"`
	StatusTransaksi TransactionStatus `gorm:"column:status_transaksi;type:transaction_status;default:Pending" json:"status_transaksi"`
	CreatedAt       time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	
	// Relasi
	Barang          Item              `gorm:"foreignKey:BarangID" json:"barang,omitempty"`
	Pembeli         User              `gorm:"foreignKey:PembeliID" json:"pembeli,omitempty"`
}

// TableName mengatur nama tabel di database
func (Transaction) TableName() string {
	return "transaksi"
}

// TransactionResponse adalah format respons untuk data transaksi
type TransactionResponse struct {
	ID              uint              `json:"id"`
	BarangID        uint              `json:"barang_id"`
	PembeliID       uint              `json:"pembeli_id"`
	TanggalTransaksi time.Time         `json:"tanggal_transaksi"`
	StatusTransaksi TransactionStatus `json:"status_transaksi"`
	CreatedAt       time.Time         `json:"created_at"`
	Barang          ItemResponse      `json:"barang,omitempty"`
	Pembeli         UserResponse      `json:"pembeli,omitempty"`
}

// ToResponse mengubah Transaction ke TransactionResponse
func (t *Transaction) ToResponse(includeBarang, includePembeli bool) TransactionResponse {
	response := TransactionResponse{
		ID:              t.ID,
		BarangID:        t.BarangID,
		PembeliID:       t.PembeliID,
		TanggalTransaksi: t.TanggalTransaksi,
		StatusTransaksi: t.StatusTransaksi,
		CreatedAt:       t.CreatedAt,
	}

	if includeBarang {
		response.Barang = t.Barang.ToResponse(false)
	}

	if includePembeli {
		response.Pembeli = t.Pembeli.ToResponse()
	}

	return response
}
package domain

import (
	"time"
)

// Chat merepresentasikan pesan chat antara penjual dan pembeli
type Chat struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PengirimID uint      `gorm:"column:pengirim_id;not null" json:"pengirim_id"`
	PenerimaID uint      `gorm:"column:penerima_id;not null" json:"penerima_id"`
	BarangID   uint      `gorm:"column:barang_id;not null" json:"barang_id"`
	Pesan      string    `gorm:"type:text;not null" json:"pesan" validate:"required"`
	Timestamp  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"timestamp"`
	Dibaca     bool      `gorm:"default:false" json:"dibaca"`
	
	// Relasi
	Pengirim   User      `gorm:"foreignKey:PengirimID" json:"pengirim,omitempty"`
	Penerima   User      `gorm:"foreignKey:PenerimaID" json:"penerima,omitempty"`
	Barang     Item      `gorm:"foreignKey:BarangID" json:"barang,omitempty"`
}

// TableName mengatur nama tabel di database
func (Chat) TableName() string {
	return "chat"
}

// ChatResponse adalah format respons untuk data chat
type ChatResponse struct {
	ID         uint          `json:"id"`
	PengirimID uint          `json:"pengirim_id"`
	PenerimaID uint          `json:"penerima_id"`
	BarangID   uint          `json:"barang_id"`
	Pesan      string        `json:"pesan"`
	Timestamp  time.Time     `json:"timestamp"`
	Dibaca     bool          `json:"dibaca"`
	Pengirim   UserResponse  `json:"pengirim,omitempty"`
	Penerima   UserResponse  `json:"penerima,omitempty"`
	Barang     ItemResponse  `json:"barang,omitempty"`
}

// ToResponse mengubah Chat ke ChatResponse
func (c *Chat) ToResponse(includePengirim, includePenerima, includeBarang bool) ChatResponse {
	response := ChatResponse{
		ID:         c.ID,
		PengirimID: c.PengirimID,
		PenerimaID: c.PenerimaID,
		BarangID:   c.BarangID,
		Pesan:      c.Pesan,
		Timestamp:  c.Timestamp,
		Dibaca:     c.Dibaca,
	}

	if includePengirim {
		response.Pengirim = c.Pengirim.ToResponse()
	}

	if includePenerima {
		response.Penerima = c.Penerima.ToResponse()
	}

	if includeBarang {
		response.Barang = c.Barang.ToResponse(false)
	}

	return response
}
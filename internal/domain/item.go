package domain

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// Status barang
type ItemStatus string

const (
	StatusTersedia ItemStatus = "Tersedia"
	StatusTerjual  ItemStatus = "Terjual"
	StatusDihapus  ItemStatus = "Dihapus"
)

// Kategori barang
type ItemCategory string

const (
	CategoryBuku      ItemCategory = "Buku"
	CategoryElektronik ItemCategory = "Elektronik"
	CategoryPerabotan  ItemCategory = "Perabotan"
	CategoryKosKosan   ItemCategory = "Kos-kosan"
	CategoryLainnya    ItemCategory = "Lainnya"
)

// Item merepresentasikan barang yang dijual
type Item struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PenjualID  uint           `gorm:"column:penjual_id;not null" json:"penjual_id"`
	NamaBarang string         `gorm:"column:nama_barang;size:100;not null" json:"nama_barang" validate:"required"`
	Harga      float64        `gorm:"type:decimal(10,2);not null" json:"harga" validate:"required,gt=0"`
	Kategori   ItemCategory   `gorm:"type:item_category;not null" json:"kategori" validate:"required,oneof=Buku Elektronik Perabotan Kos-kosan Lainnya"`
	Deskripsi  string         `gorm:"type:text" json:"deskripsi"`
	Gambar     string         `gorm:"size:255" json:"gambar"`
	Status     ItemStatus     `gorm:"type:item_status;default:Tersedia" json:"status"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-"`
	Penjual    User           `gorm:"foreignKey:PenjualID" json:"penjual,omitempty"`
}

// TableName mengatur nama tabel di database
func (Item) TableName() string {
	return "barang"
}

// ItemResponse adalah format respons untuk data barang
type ItemResponse struct {
	ID         uint         `json:"id"`
	PenjualID  uint         `json:"penjual_id"`
	NamaBarang string       `json:"nama_barang"`
	Harga      float64      `json:"harga"`
	Kategori   ItemCategory `json:"kategori"`
	Deskripsi  string       `json:"deskripsi"`
	Gambar     string       `json:"gambar"`
	Status     ItemStatus   `json:"status"`
	CreatedAt  string       `json:"created_at"`
	Penjual    *UserResponse `json:"penjual,omitempty"`
}

// ToResponse mengkonversi model Item ke respons API
func (i *Item) ToResponse(withPenjual bool) ItemResponse {
	var penjualResponse *UserResponse
	if withPenjual {
		resp := i.Penjual.ToResponse()
		penjualResponse = &resp
	}

	// Format gambar URL
	gambarURL := i.Gambar
	if gambarURL != "" {
		// Convert /download to /view if needed
		gambarURL = strings.Replace(gambarURL, "/download", "/view", 1)
	}

	return ItemResponse{
		ID:         i.ID,
		NamaBarang: i.NamaBarang,
		Harga:      i.Harga,
		Kategori:   i.Kategori,
		Deskripsi:  i.Deskripsi,
		Status:     i.Status,
		Gambar:     gambarURL,
		PenjualID:  i.PenjualID,
		Penjual:    penjualResponse,
		CreatedAt:  i.CreatedAt.Format(time.RFC3339),
	}
}
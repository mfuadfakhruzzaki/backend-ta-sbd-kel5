package domain

import (
	"fmt"
	"os"
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
	CreatedAt  time.Time    `json:"created_at"`
	Penjual    UserResponse `json:"penjual,omitempty"`
}

// ToResponse mengubah Item ke ItemResponse
func (i *Item) ToResponse(includePenjual bool) ItemResponse {
	var gambarURL string
	
	// Fix URL gambar jika perlu (pastikan URL memiliki parameter project)
	if i.Gambar != "" {
		gambarURL = i.Gambar
		// Jika URL gambar berasal dari Appwrite tapi belum memiliki parameter project
		if strings.Contains(gambarURL, "/storage/buckets/") && !strings.Contains(gambarURL, "?project=") {
			// Cari projectID dari config
			if strings.Contains(gambarURL, "/files/") {
				parts := strings.Split(gambarURL, "/files/")
				if len(parts) > 1 {
					fileIDParts := strings.Split(parts[1], "/")
					if len(fileIDParts) > 0 {
						// Use environment variable or default project ID
						projectID := os.Getenv("APPWRITE_PROJECT_ID")
						if projectID == "" {
							projectID = "67e7bbfb003b2a88a380" // Default project ID
						}
						
						// Convert /view to /download if needed
						gambarURL = strings.Replace(gambarURL, "/view", "/download", 1)
						
						// Add project parameter
						gambarURL = fmt.Sprintf("%s?project=%s", gambarURL, projectID)
					}
				}
			}
		}
	}
	
	response := ItemResponse{
		ID:         i.ID,
		PenjualID:  i.PenjualID,
		NamaBarang: i.NamaBarang,
		Harga:      i.Harga,
		Kategori:   i.Kategori,
		Deskripsi:  i.Deskripsi,
		Gambar:     gambarURL,
		Status:     i.Status,
		CreatedAt:  i.CreatedAt,
	}

	if includePenjual {
		response.Penjual = i.Penjual.ToResponse()
	}

	return response
}
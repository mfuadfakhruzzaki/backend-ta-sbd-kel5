package domain

// CreateItemRequest model untuk keperluan dokumentasi Swagger
type CreateItemRequest struct {
	NamaBarang string       `json:"nama_barang" example:"Laptop Macbook Pro 2019"`
	Harga      float64      `json:"harga" example:"10000000"`
	Kategori   ItemCategory `json:"kategori" example:"Elektronik"`
	Deskripsi  string       `json:"deskripsi" example:"Laptop dalam kondisi baik"`
}

// UpdateItemRequest model untuk keperluan dokumentasi Swagger
type UpdateItemRequest struct {
	NamaBarang string       `json:"nama_barang,omitempty" example:"Laptop Macbook Pro 2019 M1"`
	Harga      float64      `json:"harga,omitempty" example:"9500000"`
	Kategori   ItemCategory `json:"kategori,omitempty" example:"Elektronik"`
	Deskripsi  string       `json:"deskripsi,omitempty" example:"Laptop dalam kondisi baik, baru dipakai 6 bulan"`
}

// UpdateItemStatusRequest model untuk keperluan dokumentasi Swagger
type UpdateItemStatusRequest struct {
	Status ItemStatus `json:"status" example:"Terjual"`
}

// UploadImageResponse model untuk keperluan dokumentasi Swagger
type UploadImageResponse struct {
	FileName string `json:"file_name" example:"item_1_1620000000.jpg"`
	FileID   string `json:"file_id" example:"item_1_1620000000"`
	ViewURL  string `json:"view_url" example:"http://endpoint.com/storage/buckets/bucket-id/files/file-id/download?project=project-id"`
} 
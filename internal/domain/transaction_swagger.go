package domain

// CreateTransactionRequest model untuk keperluan dokumentasi Swagger
type CreateTransactionRequest struct {
	BarangID uint `json:"barang_id" example:"1" binding:"required"`
}

// UpdateTransactionStatusRequest model untuk keperluan dokumentasi Swagger
type UpdateTransactionStatusRequest struct {
	Status TransactionStatus `json:"status" example:"Selesai" binding:"required,oneof=Pending Selesai Dibatalkan"`
}

// TransactionResponse model for Swagger documentation
type TransactionResponseSwagger struct {
	ID              uint              `json:"id" example:"1"`
	BarangID        uint              `json:"barang_id" example:"5"`
	PembeliID       uint              `json:"pembeli_id" example:"2"`
	TanggalTransaksi string           `json:"tanggal_transaksi" example:"2023-05-15T14:30:45Z"`
	StatusTransaksi string            `json:"status_transaksi" example:"Pending"`
	CreatedAt       string            `json:"created_at" example:"2023-05-15T14:30:45Z"`
	Barang          ItemResponseSwagger     `json:"barang,omitempty"`
	Pembeli         UserResponseSwagger     `json:"pembeli,omitempty"`
}

// ItemResponseSwagger for embedding in TransactionResponse
type ItemResponseSwagger struct {
	ID          uint   `json:"id" example:"5"`
	Nama        string `json:"nama" example:"Laptop Asus ROG"`
	Deskripsi   string `json:"deskripsi" example:"Laptop gaming dengan spesifikasi tinggi"`
	Harga       int    `json:"harga" example:"15000000"`
	Kategori    string `json:"kategori" example:"Elektronik"`
	Status      string `json:"status" example:"Tersedia"`
	ViewURL     string `json:"view_url" example:"http://endpoint.com/storage/buckets/bucket-id/files/file-id/view?project=project-id"`
	PenjualID   uint   `json:"penjual_id" example:"3"`
	Penjual     *UserResponseSwagger `json:"penjual,omitempty"`
}

// UserResponseSwagger for embedding in TransactionResponse
type UserResponseSwagger struct {
	ID     uint   `json:"id" example:"2"`
	Email  string `json:"email" example:"user@example.com"`
	Nama   string `json:"nama" example:"John Doe"`
	NoHP   string `json:"no_hp" example:"081234567890"`
	Alamat string `json:"alamat" example:"Jl. Contoh No. 123, Jakarta"`
	Role   string `json:"role" example:"user"`
} 
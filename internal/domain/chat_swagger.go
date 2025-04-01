package domain

// SendMessageRequestSwagger model for Swagger documentation
type SendMessageRequestSwagger struct {
	PenerimaID uint   `json:"penerima_id" example:"2" binding:"required"`
	BarangID   uint   `json:"barang_id" example:"5" binding:"required"`
	Pesan      string `json:"pesan" example:"Apakah barang masih tersedia?" binding:"required"`
}

// ChatResponseSwagger model for Swagger documentation
type ChatResponseSwagger struct {
	ID         uint   `json:"id" example:"1"`
	PengirimID uint   `json:"pengirim_id" example:"1"`
	PenerimaID uint   `json:"penerima_id" example:"2"`
	BarangID   uint   `json:"barang_id" example:"5"`
	Pesan      string `json:"pesan" example:"Apakah barang masih tersedia?"`
	Dibaca     bool   `json:"dibaca" example:"false"`
	CreatedAt  string `json:"created_at" example:"2023-05-15T14:30:45Z"`
	Pengirim   UserResponseSwagger `json:"pengirim,omitempty"`
	Penerima   UserResponseSwagger `json:"penerima,omitempty"`
	Barang     ItemResponseSwagger `json:"barang,omitempty"`
}

// ChatPartnerResponseSwagger model for Swagger documentation
type ChatPartnerResponseSwagger struct {
	PartnerID   uint   `json:"partner_id" example:"2"`
	PartnerNama string `json:"partner_nama" example:"John Doe"`
	LastMessage string `json:"last_message" example:"Apakah barang masih tersedia?"`
	LastTime    string `json:"last_time" example:"2023-05-15T14:30:45Z"`
	Unread      int    `json:"unread" example:"3"`
} 
package domain

// UpdateUserRequest model untuk keperluan dokumentasi Swagger
type UpdateUserRequest struct {
	Nama     string `json:"nama,omitempty" example:"John Doe"`
	NoHP     string `json:"no_hp,omitempty" example:"08123456789"`
	Alamat   string `json:"alamat,omitempty" example:"Jl. Contoh No. 123"`
	Password string `json:"password,omitempty" example:"newpassword123"`
} 
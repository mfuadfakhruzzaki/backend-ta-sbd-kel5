package domain

// RegisterRequest model untuk keperluan dokumentasi Swagger
type RegisterRequest struct {
	Nama     string `json:"nama" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
	NoHP     string `json:"no_hp" example:"08123456789"`
	Alamat   string `json:"alamat" example:"Jl. Contoh No. 123"`
}

// RegisterResponse model untuk keperluan dokumentasi Swagger
type RegisterResponse struct {
	UserID uint `json:"user_id" example:"1"`
}

// LoginRequest model untuk keperluan dokumentasi Swagger
type LoginRequest struct {
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"password123"`
}

// UserPublic model untuk keperluan dokumentasi Swagger
type UserPublic struct {
	ID      uint   `json:"id" example:"1"`
	Nama    string `json:"nama" example:"John Doe"`
	Email   string `json:"email" example:"john@example.com"`
	NoHP    string `json:"no_hp" example:"08123456789"`
	Alamat  string `json:"alamat" example:"Jl. Contoh No. 123"`
	Role    Role   `json:"role" example:"user"`
	IsActive bool  `json:"is_active" example:"true"`
}

// LoginResponse model untuk keperluan dokumentasi Swagger
type LoginResponse struct {
	Token string     `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserPublic `json:"user"`
} 
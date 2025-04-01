package domain

// LoginRequestSwagger model for Swagger documentation
type LoginRequestSwagger struct {
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
	Password string `json:"password" example:"password123" binding:"required,min=8"`
}

// LoginResponseSwagger model for Swagger documentation
type LoginResponseSwagger struct {
	Token string            `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserResponseSwagger `json:"user"`
}

// RegisterRequestSwagger model for Swagger documentation
type RegisterRequestSwagger struct {
	Nama     string `json:"nama" example:"John Doe" binding:"required"`
	Email    string `json:"email" example:"user@example.com" binding:"required,email"`
	Password string `json:"password" example:"password123" binding:"required,min=8"`
	NoHP     string `json:"no_hp" example:"081234567890" binding:"required"`
	Alamat   string `json:"alamat" example:"Jl. Contoh No. 123, Jakarta"`
}

// RegisterResponseSwagger model for Swagger documentation
type RegisterResponseSwagger struct {
	UserID uint `json:"user_id" example:"1"`
}

// RefreshTokenRequestSwagger model for Swagger documentation
type RefreshTokenRequestSwagger struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." binding:"required"`
}

// RefreshTokenResponseSwagger model for Swagger documentation
type RefreshTokenResponseSwagger struct {
	Token        string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ResetPasswordRequestSwagger model for Swagger documentation
type ResetPasswordRequestSwagger struct {
	Email string `json:"email" example:"user@example.com" binding:"required,email"`
} 
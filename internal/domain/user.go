package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role pengguna
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User merepresentasikan pengguna dalam sistem
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Nama      string         `gorm:"size:100;not null" json:"nama" validate:"required,min=3,max=100"`
	Email     string         `gorm:"size:100;not null;uniqueIndex" json:"email" validate:"required,email"`
	Password  string         `gorm:"size:255;not null" json:"password" validate:"required,min=8"` // Pastikan tag json adalah "password" bukan "Password"
	NoHP      string         `gorm:"column:no_hp;size:15;not null" json:"no_hp" validate:"required"`
	Alamat    string         `gorm:"type:text" json:"alamat"`
	Role      Role           `gorm:"type:user_role;default:user" json:"role"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// TableName mengatur nama tabel di database
func (User) TableName() string {
	return "pengguna"
}

// HashPassword mengenkripsi password user dengan bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword memeriksa apakah password yang diberikan cocok dengan hash yang tersimpan
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Untuk keamanan, hapus password saat mengembalikan data ke client
type UserResponse struct {
	ID        uint      `json:"id"`
	Nama      string    `json:"nama"`
	Email     string    `json:"email"`
	NoHP      string    `json:"no_hp"`
	Alamat    string    `json:"alamat"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse mengubah User ke UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Nama:      u.Nama,
		Email:     u.Email,
		NoHP:      u.NoHP,
		Alamat:    u.Alamat,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
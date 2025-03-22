FROM golang:1.24-alpine AS builder

WORKDIR /app

# Salin go.mod dan go.sum terlebih dahulu untuk memanfaatkan caching
COPY go.mod go.sum ./
RUN go mod download

# Salin seluruh kode sumber
COPY . .

# Kompilasi aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# Tahap kedua: Buat image yang lebih kecil
FROM alpine:latest

WORKDIR /app

# Install paket yang diperlukan
RUN apk --no-cache add ca-certificates tzdata

# Salin binary dari tahap builder
COPY --from=builder /app/main .

# Salin file konfigurasi
COPY --from=builder /app/.env* ./

# Buat direktori uploads
RUN mkdir -p /app/uploads

# Ekspose port yang digunakan
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
# Tahap 1: Build aplikasi Golang
FROM golang:1.25.5-alpine AS builder
WORKDIR /app

# Copy file dependency dan download
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh kode dan build menjadi file binary bernama 'main'
COPY . .
RUN go build -o main .

# Tahap 2: Buat image yang sangat kecil untuk menjalankan aplikasinya
FROM alpine:3.19
WORKDIR /app

# Ambil file binary dari Tahap 1
COPY --from=builder /app/main .

# Ekspos port API kita
EXPOSE 8080

# Perintah untuk menjalankan aplikasi
CMD ["./main"]

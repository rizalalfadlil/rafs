# --- Tahap 1: Build Frontend (Vue App) ---
FROM node:22-alpine AS frontend-builder
WORKDIR /app/admin
COPY admin/package.json ./
RUN npm install
COPY admin/ ./
RUN npm run build

# --- Tahap 2: Build Backend (Go App) ---
FROM golang:1.26-alpine AS backend-builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o server-app .

# --- Tahap 3: Produksi ---
FROM alpine:latest
RUN apk add --no-cache git
WORKDIR /app

# Menyalin file biner dari tahap backend-builder
COPY --from=backend-builder /app/server-app .

# Menyalin berkas admin hasil build dari tahap frontend-builder
COPY --from=frontend-builder /app/admin/dist ./admin/dist

# Mengekspos port 8080 agar bisa diakses dari luar container
EXPOSE 8080

# Menjalankan aplikasi
CMD ["./server-app"]
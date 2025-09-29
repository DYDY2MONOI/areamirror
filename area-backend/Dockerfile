FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/uploads ./uploads

# Créer un fichier .env avec les variables d'environnement
RUN echo "DB_HOST=database" > .env && \
    echo "DB_PORT=5432" >> .env && \
    echo "DB_USER=area_user" >> .env && \
    echo "DB_PASSWORD=area_password" >> .env && \
    echo "DB_NAME=area_db" >> .env && \
    echo "JWT_SECRET=your_jwt_secret_key_here_change_in_production" >> .env

EXPOSE 8080

CMD ["./main"]

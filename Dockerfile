# Этап сборки
FROM golang:1.25 AS builder

WORKDIR /app

# Копируем go.mod и go.sum, качаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарь
# * CGO_ENABLED=0 компилято выключает исользование С, и Go собирает чистый статический бинарник. Если чистое CLI, для GUI может все сломать
# * -o указывает название выходного бинарника
RUN CGO_ENABLED=0 go build -o vado-go ./cmd/cli

# Минимальный финальный образ
FROM debian:bookworm-slim

# устанавливаем клиент и утилиту pg_isready (для wait-script)
RUN apt-get update && apt-get install -y postgresql-client ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
# копируем бинарь из билдера
COPY --from=builder /app/vado-go .
# копируем скрипт ожидания (привожу ниже)
COPY wait-for-postgres.sh .
RUN chmod +x wait-for-postgres.sh

# сперва ждём БД, потом запускаем бинарь
ENTRYPOINT ["./wait-for-postgres.sh"]
CMD ["./vado-go"]
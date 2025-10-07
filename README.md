# Build GUI
## Простая сборка исполняемого файла
```shell
CGO_ENABLED=1 GOOS=linux go build -installsuffix cgo -o vado-gui ./cmd/gui
```
1. `CGO_ENABLED=0`
* Отключает использование  **CGO** (C buildings).
* То есть компилятор **Go** не будет пытаться линковать код на С.
* Это делает бинарник **чисто статически собранным** - без зависимостей от системы С-библиотек.
* Особенно важно при сборке под **Linux**, чтобы бинарник работал на минимальных контейнерах вроде `alpine` или `scratch`.

💡 Итого: бинарь становится portable — можно запустить на любом Linux, без libc.

2. `GOOS=linux`
* Указывает целевую операционную систему для сборки.
* Позволяет собирать кроссплатформенно.
Например, если ты на macOS или Windows, то всё равно получишь бинарь под Linux.
3. `-installsuffix cgo`
* Добавляет суффикс `cgo` к пути установки (build cache, pkg dir).
* Это нужно, чтобы Go не путал обычные и `CGO_DISABLED` сборки.

Например, Go кэширует сборки пакетов.
Если ты собрал с CGO_ENABLED=1, а потом с CGO_ENABLED=0,
без -installsuffix cgo Go может взять старые артефакты из кэша.
Этот флаг гарантирует, что будет использован чистый кэш для версии без cgo.

# Linux
```shell
sudo lsof -i :9092
sudo kill -9 <PID>
```

# gRPC
Генерация из .proto файла классов go
1. Задачи
```shell
protoc -I=proto \
  --go_out=./internal/pb/taskpb \
  --go_opt=paths=source_relative \
  --go-grpc_out=./internal/pb/taskpb \
  --go-grpc_opt=paths=source_relative \
  proto/task.proto
```
2. Пользователи
```shell
   protoc --go_out=./ --go-grpc_out=./ proto/user.proto
```

# Kafka
Прочитать сообщения
```shell
docker exec -it kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic tasks \
  --from-beginning
```

# Docker compose

Запуск в фоновом режиме (`-d`) и `--build` пересоберет образ приложения (если код менялся):
```shell
docker compose up --build -d
```
Остановка (`-v` удаляет Контейнеры и Сеть и Тома, без -v тома остаются, данные БД сохраняются):
```shell
docker compose down -v
```

Проверить статус (-a показать все):
```shell
docker compose ps -a
```
Логи приложения:
```bash
docker compose logs -f app
```
Логи DB:
```bash
docker compose logs -f db
```

# Docker
Файл `Docker`
```Dockerfile
# Этап сборки
FROM golang:1.25 AS builder

WORKDIR /app

# Копируем go.mod и go.sum, качаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарь
RUN go build -o vado-go ./cmd/cli

# Минимальный финальный образ
FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/vado-go .

CMD ["./vado-go"]
```
Собираем образ
```bash
docker build . -t vado
```
Запускаем контейнер
```bash
docker run -d -p 5555:5555 vado
```
Заходим на `http://localhost:5555/`
# Golang
**go run -race main.go** определение гонки данных

**go get *link*** добавление новой зависимости

**go mod tidy** оптимизация зависимостей

**go list -m [-u] all** просмотр зависимостей текущего проекта (-u показывает доступные обновления)

**go get -u ./... && go mod tidy** обновит все зависимости проекта до последних совместимых минорных и патч-версий + почистит go.mod и go.sum от неиспользуемых пакетов.

**go list -m all | grep fyne** версии подключенных библиотек

# Godoc
**go install golang.org/x/tools/cmd/godoc@latest** установка

**godoc -http=:6060** запуск на порту (http://localhost:6060/pkg/vado/)

**wget -r -np -nH --cut-dirs=1 -P docs http://localhost:6060/pkg/vado/** генерация статической страницы.

# Idea
**Ctrl + Alt + L** Переформатировать код

# Swagger
`swag init -g ./internal/transport/rest/taskHandler.go -o internal/docs` Генерация

# Stack
* Zap logger
* Fyne
* Swagger
* Redpanda for Kafka UI
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
* zap logger
* fyne
* Swagger
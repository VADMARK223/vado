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
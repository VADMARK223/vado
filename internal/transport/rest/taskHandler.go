package rest

// @title           Vado API
// @version         1.0
// @description     Это REST API для задач.
// @host            localhost:5555
// @BasePath        /

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vado/internal/model"
	"vado/internal/service"
	"vado/pkg/logger"

	"github.com/k0kubun/pp"
	"go.uber.org/zap"
)

const (
	slowRequestDelaySecond = 10 // Длительность выполнения медленного запроса// Время "мягкой" остановки сервер
)

type TaskHandler struct {
	Service service.ITaskService
}

func (th *TaskHandler) TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAllTasks(th, w, r)
	case http.MethodPost:
		createTask(th, w, r)
	case http.MethodDelete:
		deleteAllTasks(th, w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

// @Summary      Удалить все задачи
// @Description  Чистит список задач
// @Tags         tasks
// @Produce      plain
// @Success      200  {string}  string  "Delete tasks."
// @Router       /tasks [delete]
func deleteAllTasks(th *TaskHandler, w http.ResponseWriter, _ *http.Request) {
	th.Service.DeleteAllTasks()
	_, _ = w.Write([]byte("Delete tasks."))
}

// @Summary      Создать задачу
// @Description  Добавляет новую задачу
// @Tags         tasks
// @Accept       json
// @Produce      plain
// @Param        task  body      model.Task  true  "Task info"
// @Success      200   {string}  string      "Success."
// @Router       /tasks [post]
func createTask(th *TaskHandler, w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if task.Name == "" {
		logger.L().Error("create task failed")
		http.Error(w, "Task name empty.", http.StatusBadRequest)
		return
	}

	err := th.Service.CreateTask(task)
	if err != nil {
		logger.L().Error("create task failed", zap.Error(err))
		http.Error(w, "failed to create task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = pp.Println("Created task", task)
	_, err = w.Write([]byte("Success."))
	if err != nil {
		http.Error(w, "failed to write: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary      Получить список задач.
// @Description  Возвращает все задачи
// @Tags         tasks
// @Produce      json
// @Success      200  {array}  model.Task
// @Router       /tasks [get]
func getAllTasks(th *TaskHandler, w http.ResponseWriter, _ *http.Request) {
	tasksList, err := th.Service.GetAllTasks()
	if err != nil {
		http.Error(w, "failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasksList); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (th *TaskHandler) TaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Обрезаем "/tasks/"
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid id=%s.", idStr), http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getTaskByID(th, id, w)
	case http.MethodDelete:
		deleteTaskByID(th, id, w)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

// @Summary      Удалить задачу по ID
// @Description  Удаляет задачу по её ID
// @Tags         tasks
// @Param        id   path      int  true  "ID задачи"
// @Success      200  {string}  string  "Delete task {id}."
// @Failure      400  {string}  string  "Invalid id"
// @Failure      404  {string}  string  "Task not found"
// @Router       /tasks/{id} [delete]
func deleteTaskByID(th *TaskHandler, id int, w http.ResponseWriter) {
	err := th.Service.DeleteTask(id)
	if err != nil {
		logger.L().Warn("Failed to delete task:", zap.Error(err))
		http.Error(w, "Failed to delete task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	msg := fmt.Sprintf("Delete task %d.", id)
	_, _ = w.Write([]byte(msg))
	logger.L().Info(msg)
}

// @Summary      Получить задачу по ID
// @Description  Возвращает задачу по её ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "ID задачи"
// @Success      200  {object}  model.Task
// @Failure      400  {string}  string  "Invalid id"
// @Failure      404  {string}  string  "Task not found"
// @Router       /tasks/{id} [get]
func getTaskByID(th *TaskHandler, taskId int, w http.ResponseWriter) {
	task, err := th.Service.GetTaskByID(taskId)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (th *TaskHandler) SlowHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Started slow request...")
	time.Sleep(time.Second * slowRequestDelaySecond)
	str := "Hello from slow handler!"
	_, err := w.Write([]byte(str))
	if err != nil {
		fmt.Println("Error", err)
	} else {
		fmt.Println("Finished slow request")
	}
}

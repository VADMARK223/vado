package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"vado/internal/model"
	"vado/internal/service"

	"github.com/k0kubun/pp"
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
	case http.MethodPost:
		var task model.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
			return
		}

		_, _ = pp.Println("Create task", task)

		if task.Name == "" {
			http.Error(w, "Task name empty.", http.StatusBadRequest)
			return
		}

		err := th.Service.CreateTask(task)
		if err != nil {
			http.Error(w, "failed to create task: "+err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write([]byte("Task created."))
		if err != nil {
			http.Error(w, "failed to write: "+err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		th.Service.DeleteAllTasks()
		_, _ = w.Write([]byte("Delete tasks."))
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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
		_, _ = w.Write([]byte("Not implement."))
	case http.MethodDelete:
		err := th.Service.DeleteTask(id)
		if err != nil {
			http.Error(w, "Failed to delete task: "+err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(fmt.Sprintf("Delete task %d.", id)))
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
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

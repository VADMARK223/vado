package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"vado/internal/service"
)

const (
	slowRequestDelaySecond = 10 // Длительность выполнения медленного запроса// Время "мягкой" остановки сервер
)

type TaskHandler struct {
	Service service.ITaskService
}

func (th *TaskHandler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		msg := "Error: Method not allowed."
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Println(msg)
		if _, err := w.Write([]byte(msg)); err != nil {
			fmt.Println("Error writing response:", err)
		}
		return
	}

	tasksList, err := th.Service.GetAllTasks()
	if err != nil {
		http.Error(w, "failed to get tasks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasksList); err != nil {
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

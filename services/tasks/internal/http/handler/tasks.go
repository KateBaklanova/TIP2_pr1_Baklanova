package handler

import (
	"encoding/json"
	"kate/services/tasks/internal/service"
	"net/http"
	"strings"
)

type TaskHandler struct {
	taskSvc *service.TaskService
}

func NewTaskHandler(ts *service.TaskService) *TaskHandler {
	return &TaskHandler{taskSvc: ts}
}

func (h *TaskHandler) handleError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var task service.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	created := h.taskSvc.Create(task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	tasks := h.taskSvc.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	if id == "" {
		h.handleError(w, http.StatusBadRequest, "missing id")
		return
	}

	task, ok := h.taskSvc.GetByID(id)
	if !ok {
		h.handleError(w, http.StatusNotFound, "task not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		h.handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	if id == "" {
		h.handleError(w, http.StatusBadRequest, "missing id")
		return
	}

	var updates service.Task
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		h.handleError(w, http.StatusBadRequest, "invalid json")
		return
	}

	updated, ok := h.taskSvc.Update(id, updates)
	if !ok {
		h.handleError(w, http.StatusNotFound, "task not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.handleError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	if id == "" {
		h.handleError(w, http.StatusBadRequest, "missing id")
		return
	}

	ok := h.taskSvc.Delete(id)
	if !ok {
		h.handleError(w, http.StatusNotFound, "task not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

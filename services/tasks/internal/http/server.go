package http

import (
	"kate/services/tasks/internal/client"
	"kate/services/tasks/internal/http/handler"
	"kate/services/tasks/internal/service"
	"kate/shared/middleware"
	"log"
	"net/http"
	"time"
)

func StartServer(port string, authBaseURL string) {
	taskSvc := service.NewTaskService()
	authCli := client.NewAuthClient(authBaseURL, 3*time.Second)
	taskHandler := handler.NewTaskHandler(taskSvc)

	mux := http.NewServeMux()

	// Публичные маршруты (если есть)
	// Защищенные маршруты с проверкой токена через middleware
	mux.Handle("/v1/tasks", middleware.RequestIDMiddleware(
		handler.AuthMiddleware(authCli)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					taskHandler.GetAllTasks(w, r)
				case http.MethodPost:
					taskHandler.CreateTask(w, r)
				default:
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				}
			}),
		),
	))

	mux.Handle("/v1/tasks/", middleware.RequestIDMiddleware(
		handler.AuthMiddleware(authCli)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					taskHandler.GetTaskByID(w, r)
				case http.MethodPatch:
					taskHandler.UpdateTask(w, r)
				case http.MethodDelete:
					taskHandler.DeleteTask(w, r)
				default:
					http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				}
			}),
		),
	))

	// Логирование для всех запросов
	handlerWithLogging := middleware.LoggingMiddleware(mux)

	log.Printf("Tasks service starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handlerWithLogging); err != nil {
		log.Fatal(err)
	}
}

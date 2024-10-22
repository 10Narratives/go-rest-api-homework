package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// getTasks handles the HTTP request to retrieve a list of tasks.
// It writes the tasks in JSON format to the provided http.ResponseWriter.
//
// In case of an error during the JSON marshaling process,
// it responds with an HTTP 500 Internal Server Error and
// writes the error message to the response body.
//
// Parameters:
//   - writer: The http.ResponseWriter used to write the response.
//   - req: The http.Request received from the client. This parameter is ignored in this function.
func getTasks(writer http.ResponseWriter, _ *http.Request) {
	response, err := json.Marshal(tasks)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

// getTask retrieves a task by its ID from the provided URL parameter
// and writes its JSON representation to the http.ResponseWriter. If
// the task is not found, it sends a HTTP 400 Bad Request response.
// In case of an error during JSON marshaling, it responds with
// a HTTP 400 Bad Request along with the error message.
//
// Parameters:
//   - writer: The http.ResponseWriter used to send the response to the client.
//   - request: The http.Request object that contains the URL parameters,
//     including the task ID to be retrieved.
func getTask(writer http.ResponseWriter, request *http.Request) {
	targetID := chi.URLParam(request, "id")
	task, wasFound := tasks[targetID]
	if !wasFound {
		http.Error(writer, "Task with given ID was not found", http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

// postTask handles the creation of a new task. It reads the task's
// information from the HTTP request body, unmarshals it into a
// Task struct, and stores it in the tasks map. Upon successful
// creation, it responds with a HTTP 201 Created status. If any
// errors occur during reading the request body or unmarshaling,
// it responds with a HTTP 400 Bad Request along with the error message.
//
// Parameters:
//   - writer: The http.ResponseWriter used to send the response to the client.
//   - request: The http.Request object that contains the new task data in the body.
func postTask(writer http.ResponseWriter, request *http.Request) {
	var newTask Task
	var buffer bytes.Buffer

	_, err := buffer.ReadFrom(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buffer.Bytes(), &newTask); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
}

// deleteTask handles the deletion of a task identified by the
// task ID provided in the URL parameter. If the task with the
// specified ID does not exist, it responds with a HTTP 400
// Bad Request status and an error message. Upon successful
// deletion, it returns a HTTP 200 OK status.
//
// Parameters:
//   - writer: The http.ResponseWriter used to send the response to the client.
//   - request: The http.Request object that contains the URL parameters, including the task ID.
func deleteTask(writer http.ResponseWriter, request *http.Request) {
	taskID := chi.URLParam(request, "id")

	_, wasFound := tasks[taskID]
	if !wasFound {
		http.Error(writer, "Task was not found.", http.StatusBadRequest)
		return
	}

	delete(tasks, taskID)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}

func main() {
	router := chi.NewRouter()

	router.Get("/tasks", getTasks)
	router.Post("/tasks", postTask)
	router.Get("/tasks/{id}", getTask)
	router.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
	fmt.Println("Listen and Serve")
}

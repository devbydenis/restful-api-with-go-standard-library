package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"rest-api/internal/taskstore"
	"rest-api/model"
	"strconv"
	"time"
)

type taskServer struct {
	store *taskstore.TaskStore
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (ts *taskServer) createHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling task create at %s\n", req.URL.Path)

	// enforce a JSON Content-Type
	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediaType != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	req.Body = http.MaxBytesReader(w, req.Body, 1048576)
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	var rt model.RequestTask
	if err := dec.Decode(&rt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := ts.store.CreateTask(rt.Text, rt.Tags)
	renderJSON(w, model.ResponseTask{
		Status:  http.StatusCreated,
		Message: fmt.Sprintf("Task with %d is created", id),
	})
}

func (ts *taskServer) getTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get task handler at %s\n", req.URL.Path)

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := ts.store.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, task)
}

func (ts *taskServer) getAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all task at %v", req.URL.Path)

	allTask := ts.store.GetAllTasks()

	renderJSON(w, model.ResponseTask{
		Status: http.StatusOK,
		Message: "Fetched successfully",
		Data: allTask,
	})
}

func (ts *taskServer) deleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete task at %v", req.URL.Path)

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ts.store.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, model.ResponseTask{
		Status: http.StatusNoContent,
		Message: fmt.Sprintf("Task with id %d success to remove", id),
	})
}

func (ts *taskServer) deleteAllTasksHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handle delete all task at %v", req.URL.Path)

	err := ts.store.DeleteAllTask()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderJSON(w, model.ResponseTask{
		Status: http.StatusNoContent,
		Message: "Successfully remove all task!",
	})
}

func (ts *taskServer) tagHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handle tag at %v", req.URL.Path)

	tag := req.PathValue("tag")
	tags := ts.store.GetTaskByTag(tag)

	renderJSON(w, model.ResponseTask{
		Status: http.StatusOK,
		Message: "Succesfully fetched data by tag",
		Data: tags,
	})
}

func (ts *taskServer) dueHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handlling task by due at %v", req.URL.Path)

	year, errYear := strconv.Atoi(req.PathValue("year"))
	month, errMonth := strconv.Atoi(req.PathValue("month"))
	day, errDay := strconv.Atoi(req.PathValue("day"))
	if errYear != nil || errMonth != nil || errDay != nil || month < int(time.January) || month > int(time.December) {
		http.Error(w, fmt.Sprintf("expect /due/<year>/<month/<day>, got %v", req.URL.Path), http.StatusBadRequest)
	}

	tasks := ts.store.GetTaskByDueDate(year, time.Month(month), day)
	renderJSON(w, model.ResponseTask{
		Status: http.StatusOK,
		Message: "Successfully fetch data by due",
		Data: tasks,
	})
}

func NewTaskServer(store *taskstore.TaskStore) *taskServer {
	return &taskServer{store: store}
}

func main() {
	mux := http.NewServeMux()
	server := NewTaskServer(taskstore.NewTaskStore())

	mux.HandleFunc("POST /task/", server.createHandler)
	mux.HandleFunc("GET /tasks/", server.getAllTasksHandler)
	mux.HandleFunc("GET /task/{id}/", server.getTaskHandler)
	mux.HandleFunc("GET /tag/{tag}/", server.tagHandler)
	mux.HandleFunc("GET /due/{year}/{month}/{day}/", server.dueHandler)
	mux.HandleFunc("DELETE /tasks/", server.deleteAllTasksHandler)
	mux.HandleFunc("DELETE /task/{id}/", server.deleteTaskHandler)

	port := os.Getenv("SERVERPORT") // karna di go ga kenal .env maka kita langsung tulis SERVERPORT=8080 go run main.go atau kalo mau pake .env bisa pake package go get github.com/joho/godotenv
	if port == "" {
		log.Fatal("SERVERPORT env is not set")
	}

	fmt.Println("Server is running")
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

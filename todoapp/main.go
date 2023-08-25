package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Todo struct {
	Id    int `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}
type Response struct {
	Status       string
	ErrorMessage string
	Data		 interface{}
}

var todos = make(map[int]Todo)

func Respond(status,errorMsg string,data interface{})[]byte{
	var response Response
	response.Status = status
	response.ErrorMessage = errorMsg
	response.Data = data
	body, marshErr := json.Marshal(response)
	if marshErr != nil {
		log.Println(marshErr)
	}
	return body
}
func SaveToDo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	
	body, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		log.Println(readErr)
	}
	unMarshErr := json.Unmarshal(body, &todo)
	if unMarshErr != nil {
		
		w.Write(Respond("FAILED",unMarshErr.Error(),nil))
	}
	todos[todo.Id] = todo
	w.WriteHeader(200)
	w.Write(Respond("SUCCESS","",todo))
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong!"))
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", Ping)
	mux.HandleFunc("/save",SaveToDo)

	err := http.ListenAndServe(":9000", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed!")
	} else if err != nil {
		log.Println(err)
	}
}

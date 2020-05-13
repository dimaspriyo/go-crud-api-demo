package main

import (
	"context"
	"encoding/json"
	"go-crud-api-demo/repository"
	"net/http"
	"strconv"
	"strings"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
)

//Response Struct
type Response struct {
	Message string              `json:"message,omitempty"`
	Data    []repository.Person `json:"data,omitempty"`
	Error   string              `json:"error,omitempty"`
}

var person repository.Person
var response Response

func main() {
	mux := http.NewServeMux()

	//v1 just plain API without authentication
	mux.HandleFunc("/v1", listV1)
	mux.HandleFunc("/v1/insert", insertV1)
	mux.HandleFunc("/v1/update/", updateV1)
	mux.HandleFunc("/v1/delete", deleteV1)

	//v1 API authentication using JWT
	mux.HandleFunc("/v2", listV2)
	mux.HandleFunc("/v2/insert", insertV2)
	mux.HandleFunc("/v2/update", updateV2)
	mux.HandleFunc("/v2/delete", deleteV2)

	http.ListenAndServe(":8080", mux)

}

func listV1(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	db, err := repository.NewConn(ctx)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	list, err := repository.List(db)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	response.Data = list
	response.Message = "Get All Person Data Success"
	jsonResponse(w, response)
}
func insertV1(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&person)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	db, err := repository.NewConn(ctx)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	res, err := repository.Insert(ctx, db, person)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	response.Message = "Insert Success"
	response.Data = []repository.Person{res}

}
func updateV1(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)

	rawID := lastURI(r.RequestURI)
	//convert string to int64
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	err = decoder.Decode(&person)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	db, err := repository.NewConn(ctx)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	res, err := repository.Update(ctx, db, person, id)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	response.Message = "Update Success"
	response.Data = []repository.Person{res}
	jsonResponse(w, response)

}
func deleteV1(w http.ResponseWriter, r *http.Request) {

}

func listV2(w http.ResponseWriter, r *http.Request) {

}
func insertV2(w http.ResponseWriter, r *http.Request) {

}
func updateV2(w http.ResponseWriter, r *http.Request) {

}
func deleteV2(w http.ResponseWriter, r *http.Request) {

}

//=========

func lastURI(RequestURI string) string {

	last := strings.Split(RequestURI, "/")

	ret := last[len(last)-1]

	return ret

}

func jsonResponse(w http.ResponseWriter, response Response) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

package main

import (
	"net/http"
)

func main() {

	//v1 just plain API without authentication
	http.HandleFunc("v1", listV1)
	http.HandleFunc("v1/insert", insertV1)
	http.HandleFunc("v1/update", updateV1)
	http.HandleFunc("v1/delete", deleteV1)

	//v1 API authentication using JWT
	http.HandleFunc("v2", listV2)
	http.HandleFunc("v2/insert", insertV2)
	http.HandleFunc("v2/update", updateV2)
	http.HandleFunc("v2/delete", deleteV2)

	http.ListenAndServe(":8080", nil)

}

func listV1(w http.ResponseWriter, r *http.Request) {

}
func insertV1(w http.ResponseWriter, r *http.Request) {

}
func updateV1(w http.ResponseWriter, r *http.Request) {

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

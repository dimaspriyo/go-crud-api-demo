package main

import (
	"context"
	"encoding/json"
	"errors"
	"go-crud-api-demo/config"
	"go-crud-api-demo/repository"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	Token   string              `json:"token,omitempty"` //for JWT Token
}

//Key Request Struct
type KeyRequest struct {
	Key string `json:"key"`
}

func main() {
	mux := http.NewServeMux()

	//v1 just plain API without authentication
	mux.HandleFunc("/v1", listV1)
	mux.HandleFunc("/v1/insert", insertV1)
	mux.HandleFunc("/v1/update/", updateV1)
	mux.HandleFunc("/v1/delete/", deleteV1)

	//v1 API authentication using JWT
	mux.HandleFunc("/v2/token", tokenV2)

	handlerListV2 := http.HandlerFunc(listV2)
	mux.Handle("/v2", JWtAuthenticationMiddleware(handlerListV2))

	handlerInsertV2 := http.HandlerFunc(insertV2)
	mux.Handle("/v2/insert", JWtAuthenticationMiddleware(handlerInsertV2))

	handlerUpdateV2 := http.HandlerFunc(updateV2)
	mux.Handle("/v2/update/", JWtAuthenticationMiddleware(handlerUpdateV2))

	handlerDeleteV2 := http.HandlerFunc(deleteV2)
	mux.Handle("/v2/delete/", JWtAuthenticationMiddleware(handlerDeleteV2))

	http.ListenAndServe(":8080", mux)

}

func listV1(w http.ResponseWriter, r *http.Request) {

	var response Response

	if r.Method != "GET" {
		response.Message = "Invalid HTTP Method"
		jsonResponse(w, response)
		return
	}

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

	var person repository.Person
	var response Response

	if r.Method != "POST" {
		response.Message = "Invalid HTTP Method"
		jsonResponse(w, response)
		return
	}

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

	var person repository.Person
	var response Response

	if r.Method != "PUT" {
		response.Message = "Invalid HTTP Method"
		jsonResponse(w, response)
		return
	}

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

	var response Response

	if r.Method != "DELETE" {
		response.Message = "Invalid HTTP Method"
		jsonResponse(w, response)
		return
	}

	ctx := r.Context()

	rawID := lastURI(r.RequestURI)
	//convert string to int64
	id, err := strconv.ParseInt(rawID, 10, 64)

	db, err := repository.NewConn(ctx)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	err = repository.Delete(ctx, db, id)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	response.Message = "Delete Success"
	jsonResponse(w, response)

}

//================================

func JWtAuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var response Response

			token := r.Header.Get("Authorization")

			err := checkJWTToken(token)
			if err != nil {
				response.Error = err.Error()
				jsonResponse(w, response)
				return
			}

			next.ServeHTTP(w, r)
		})
}

func tokenV2(w http.ResponseWriter, r *http.Request) {

	var response Response
	var keyRequest KeyRequest

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&keyRequest)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	err = checkRequestKey(keyRequest.Key)
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}

	response.Message = "Here's your Token"

	token, err := generateJWTToken()
	if err != nil {
		response.Error = err.Error()
		jsonResponse(w, response)
		return
	}
	response.Token = token

	jsonResponse(w, response)
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

func checkRequestKey(key string) (err error) {

	var JSONConfig config.JSONConfig
	JSONConfig, err = config.ReadJSONConfig()
	if err != nil {
		return err
	}

	//Check Key for requesting token
	if key != JSONConfig.Key {
		return errors.New("Request Key Incorrect")
	}

	return nil

}

func generateJWTToken() (token string, err error) {

	var JSONConfig config.JSONConfig
	JSONConfig, err = config.ReadJSONConfig()
	if err != nil {
		return token, err
	}

	generate := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": "admin"})

	token, err = generate.SignedString([]byte(JSONConfig.Secret))
	return token, nil

}

func checkJWTToken(myToken string) (err error) {

	myToken = strings.Replace(myToken, "Bearer ", "", -1)
	var JSONConfig config.JSONConfig
	JSONConfig, err = config.ReadJSONConfig()
	if err != nil {
		return err
	}

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(JSONConfig.Secret), nil
	})

	_, ok := token.Claims.(jwt.MapClaims)
	if !token.Valid || !ok {
		return errors.New("Invalid Token")
	}

	return nil

}

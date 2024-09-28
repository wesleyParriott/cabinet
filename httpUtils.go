package main

import "net/http"

// 200s
func Okay(response http.ResponseWriter, message []byte) {
	response.Write(message)
}

func Created(response http.ResponseWriter, location string) {
	response.WriteHeader(http.StatusCreated)
	response.Header().Set("Location", location)
	response.Write([]byte("Created " + location))
}

// 300s
// ..

// 400s
func BadRequest(response http.ResponseWriter) {
	http.Error(response, http.StatusText(400), 400)
}

func Forbidden(response http.ResponseWriter) {
	http.Error(response, http.StatusText(403), 403)
}

func NotFound(response http.ResponseWriter) {
	http.Error(response, http.StatusText(404), 404)
}

func Conflict(response http.ResponseWriter) {
	http.Error(response, http.StatusText(409), 409)
}

func EntityTooLarge(response http.ResponseWriter) {
	http.Error(response, http.StatusText(413), 413)
}

// 500s
func InternalError(response http.ResponseWriter) {
	http.Error(response, http.StatusText(500), 500)
}

func NotImplemented(response http.ResponseWriter) {
	http.Error(response, http.StatusText(501), 501)
}

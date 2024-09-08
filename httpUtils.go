package main

import "net/http"

// 200s
func Okay(response http.ResponseWriter, message []byte) {
	response.Write(message)
}

// 300s
// ...

// 400s
func BadRequest(response http.ResponseWriter) {
	http.Error(response, http.StatusText(400), 400)
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

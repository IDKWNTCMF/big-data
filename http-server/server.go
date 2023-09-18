package main

import (
	"io"
	"net/http"
	"os"
)

func replaceHandler(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		err = os.WriteFile(filename, body, 0644)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	}
}

func getHandler(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		stored, err := os.ReadFile(filename)
		if err != nil {
			return
		}
		_, err = w.Write(stored)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	}
}

func main() {
	filename := "storage"
	http.HandleFunc("/replace", replaceHandler(filename))
	http.HandleFunc("/get", getHandler(filename))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

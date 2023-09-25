package main

import (
	"io"
	"net/http"
)

func ReplaceHandler(manager *TransactionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		manager.PushTransaction(string(body))
		w.WriteHeader(200)
	}
}

func GetHandler(manager *TransactionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		stored := manager.GetStored()
		_, err := w.Write([]byte(stored))
		if err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	}
}

func main() {
	manager := NewTransactionManager()
	go manager.StartTransactionManager()
	http.HandleFunc("/replace", ReplaceHandler(&manager))
	http.HandleFunc("/get", GetHandler(&manager))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

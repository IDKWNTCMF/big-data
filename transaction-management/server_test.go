package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReplaceStatusCode(t *testing.T) {
	manager := NewTransactionManager()
	go manager.StartTransactionManager()
	reader := strings.NewReader("kek")
	req := httptest.NewRequest(http.MethodPost, "/replace", reader)
	w := httptest.NewRecorder()
	replace := ReplaceHandler(&manager)
	replace(w, req)
	res := w.Result()
	if res.StatusCode != 200 {
		t.Errorf("Expected status code to be 200 got %v", res.StatusCode)
	}
}

func TestSimpleQuery(t *testing.T) {
	manager := NewTransactionManager()
	go manager.StartTransactionManager()
	reader := strings.NewReader("kek")
	req := httptest.NewRequest(http.MethodPost, "/replace", reader)
	w := httptest.NewRecorder()
	replace := ReplaceHandler(&manager)
	replace(w, req)
	req = httptest.NewRequest(http.MethodGet, "/get", nil)
	w = httptest.NewRecorder()
	get := GetHandler(&manager)
	get(w, req)
	res := w.Result()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Expected error to be nil got %v", err)
	}
	if res.StatusCode != 200 {
		t.Errorf("Expected status code to be 200 got %v", res.StatusCode)
	}
	if string(data) != "kek" {
		t.Errorf("Expected kek got %v", string(data))
	}
}

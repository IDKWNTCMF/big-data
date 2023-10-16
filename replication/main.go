package main

import (
	"embed"
	"encoding/json"
	"io"
	"net/http"
	ws "nhooyr.io/websocket"
)

var Source string = "127.0.0.1"
var Port string = "8080"
var Counter uint64 = 1

//go:embed index.html
var content embed.FS

func main() {
	manager := NewTransactionManager()
	go manager.StartTransactionManager()
	initPeers();
	go websocketClient(&manager);
	http.Handle("/test/", http.StripPrefix("/test/", http.FileServer(http.FS(content))))
	http.HandleFunc("/vclock", handleVClock(&manager))
	http.HandleFunc("/replace", handleReplace(&manager))
	http.HandleFunc("/get", handleGet())
	http.HandleFunc("/ws", handleWS())
	err := http.ListenAndServe(":" + Port, nil)
	if err != nil {
		panic(err)
	}
}

func handleVClock(manager *TransactionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes, err := json.Marshal(manager.VClock)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		_, err = w.Write(bytes)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}
}

func handleReplace(manager *TransactionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		var id = Counter
		Counter++
		var transaction = Transaction{Source: Source, Id: id, Payload: string(body)}
		manager.Transactions <- transaction
	}
}

func handleGet() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte(Snapshot))
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}
}

func handleWS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := ws.Accept(w, r, &ws.AcceptOptions{
			InsecureSkipVerify: true,
			OriginPatterns:     []string{"*"},
		})
		if err != nil {
			panic(err)
		}
		ctx := r.Context()
		handlePeerDownstream(c, ctx, r.Host)
		c.Close(ws.StatusNormalClosure, "Connection closed")
	}
}
package main

import (
	"context"
	"fmt"
	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

func handlePeerDownstream(c *ws.Conn, ctx context.Context, peer string) {
	lastProcessed := 0
	for {
		for ; lastProcessed < len(Journal); lastProcessed++ {
			transaction := Journal[lastProcessed]
			err := wsjson.Write(ctx, c, transaction)
			if err != nil {
				break
			}
		}
		time.Sleep(time.Second)
	}
}

var peers []string

func initPeers() {
	peers = append(peers, "127.0.0.1:8081")
}

func websocketClient(manager *TransactionManager) {
	for _, peer := range peers {
		go handlePeerUpstream(manager, peer)
	}
}

func handlePeerUpstream(manager *TransactionManager, peer string) {
	for {
		ctx := context.Background()
		url := fmt.Sprintf("ws://%s/ws", peer)
		c, _, err := ws.Dial(ctx, url, nil)
		if err != nil {
			continue
		}
		acceptReplication(manager, c, ctx, peer)
		time.Sleep(15 * time.Second)
	}
}

func acceptReplication(manager *TransactionManager, c *ws.Conn, ctx context.Context, peer string) {
	for {
		var transaction Transaction
		err := wsjson.Read(ctx, c, &transaction)
		if err != nil {
			return
		}
		manager.Transactions <- transaction
	}
}
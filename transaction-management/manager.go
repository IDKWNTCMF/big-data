package main

import (
	"sync"
	"time"
)

var Stored string
var Journal []string
var Snapshot string

type TransactionManager struct {
	journals     [][]string
	transactions chan string
	mutex        sync.Mutex
	ticker       *time.Ticker
}

func NewTransactionManager() TransactionManager {
	return TransactionManager{transactions: make(chan string), ticker: time.NewTicker(time.Minute)}
}

func (manager *TransactionManager) StartTransactionManager() {
	for {
		select {
		case transaction := <-manager.transactions:
			manager.mutex.Lock()
			Journal = append(Journal, transaction)
			Stored = transaction
			manager.mutex.Unlock()
		case <-manager.ticker.C:
			manager.mutex.Lock()
			Snapshot = Stored
			manager.journals = append(manager.journals, Journal)
			Journal = nil
			manager.mutex.Unlock()
		}
	}
}

func (manager *TransactionManager) GetStored() string {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	return Stored
}

func (manager *TransactionManager) PushTransaction(transaction string) {
	manager.transactions <- transaction
}

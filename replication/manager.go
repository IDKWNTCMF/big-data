package main

import (
	jsonpatch "github.com/evanphx/json-patch/v5"
	"log"
	"sync"
)

var Snapshot string = "{}"
var Journal []Transaction

type Transaction struct {
	Source  string
	Id      uint64
	Payload string
}

type TransactionManager struct {
	Transactions chan Transaction
	Mutex        sync.Mutex
	VClock       map[string]uint64
}

func NewTransactionManager() TransactionManager {
	return TransactionManager{Transactions: make(chan Transaction), VClock: make(map[string]uint64)}
}

func (manager *TransactionManager) StartTransactionManager() {
	for {
		manager.ApplyTransaction()
	}
}

func (manager *TransactionManager) ApplyTransaction() {
	transaction := <-manager.Transactions
	manager.Mutex.Lock()
	defer manager.Mutex.Unlock()
	log.Printf("Got transaction: {Source: %s, Id: %v, Payload: %s}\n", transaction.Source, transaction.Id, transaction.Payload)

	transactionAlreadyApplied := manager.VClock[transaction.Source] >= transaction.Id
	
	if transactionAlreadyApplied {
		log.Printf("Transaction already applied\n")
		return
	}

	manager.VClock[transaction.Source] = transaction.Id
	Journal = append(Journal, transaction)
	
	patch, err := jsonpatch.DecodePatch([]byte(transaction.Payload))
	if err != nil {
		log.Printf("Transaction cannot be implied: %s", err)
		return
	}

	newsnap, err := patch.Apply([]byte(Snapshot))
	if err != nil {
		log.Printf("Transaction cannot be implied: %s", err)
		return
	}
	Snapshot = string(newsnap)
	log.Printf("Transaction applied\n")
}
package main

import (
	"ddas/lesson2/transaction"
	"io"
	"net/http"
	"sync"
)

// Очередь, которую будет читать менеджер транзакций
var transactionQueue chan []byte

// Последнее полученное body
var body []byte
var bodyMutex sync.Mutex

func main() {
	transactionQueue = make(chan []byte)
	transaction.NewTM(transactionQueue)

	mux := http.NewServeMux()

	mux.Handle("/get", http.HandlerFunc(getHandler))
	mux.Handle("/replace", http.HandlerFunc(replaceHandler))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	bodyMutex.Lock()
	defer bodyMutex.Unlock()

	_, err := w.Write(body)
	if err != nil {
		panic(err)
	}
}

func replaceHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	bodyMutex.Lock()
	defer bodyMutex.Unlock()

	transactionQueue <- data
	body = data
}

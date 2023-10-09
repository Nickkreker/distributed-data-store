package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Transaction struct {
	Source  string // фамилия
	Id      uint64 // возрастающий счетчик
	Payload string // транзакция
}

var (
	peers                   []string          = []string{"127.0.0.1:8081"} // соседи
	snap                    string            = "{\"0\": null}"            // снапшот
	wal                     []*Transaction    = make([]*Transaction, 0)    // журнал транзакций
	vclock                  map[string]uint64 = make(map[string]uint64)    // логические часы
	source                  string            = "Berezikov-2"
	localTId                uint64            = 0 // локальный счетчик транзакций
	transactionManagerQueue chan *Transaction = make(chan *Transaction)
	//go:embed html/index.html
	testResponse string

	wsConnections []*websocket.Conn = make([]*websocket.Conn, 0)
)

func main() {
	mux := http.NewServeMux()

	go transactionManagerJob()

	for _, peer := range peers {
		go websocketClientJob(peer)
	}

	mux.HandleFunc("/test", testHandler)
	mux.HandleFunc("/vclock", vclockHandler)
	mux.HandleFunc("/replace", replaceHandler)
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/ws", wsHandler)

	_ = http.ListenAndServe(":8080", mux)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(testResponse))
	if err != nil {
		panic(err)
	}
}

func vclockHandler(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(vclock)
	if err != nil {
		panic(err)
	}

	_, err = w.Write(res)
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

	localTId += 1
	transaction := &Transaction{
		Payload: string(data),
		Id:      localTId,
		Source:  source,
	}

	transactionManagerQueue <- transaction
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(snap))
	if err != nil {
		panic(err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	})

	if err != nil {
		panic(err)
	}

	slog.Info("Received new connection ", "from", r.RemoteAddr)
	wsConnections = append(wsConnections, c)
}

func transactionManagerJob() {
	for {
		transaction := <-transactionManagerQueue

		if vclock[transaction.Source] > transaction.Id {
			continue
		}
		vclock[transaction.Source] = transaction.Id + 1
		wal = append(wal, transaction)
		patch, err := jsonpatch.DecodePatch([]byte(transaction.Payload))
		if err != nil {
			panic(err)
		}

		snapBytes, err := patch.Apply([]byte(snap))
		if err != nil {
			panic(err)
		}

		snap = string(snapBytes)

		slog.Info("Sending transaction to peers", "transaction", *transaction)
		for _, conn := range wsConnections {
			err = wsjson.Write(context.Background(), conn, transaction)
			if err != nil {
				panic(err)
			}
		}
	}
}

func websocketClientJob(peer string) {
	var conn *websocket.Conn
	var err error
	ctx := context.Background()
	for {
		conn, _, err = websocket.Dial(ctx, fmt.Sprintf("ws://%s/ws", peer), nil)
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	for {
		var transaction Transaction
		err = wsjson.Read(ctx, conn, &transaction)
		if errors.Is(err, io.EOF) {
			slog.Info("Peer disconnected", "peer", peer)
			break
		}

		slog.Info("Received", "transaction", transaction)
		if err != nil {
			panic(err)
		}

		transactionManagerQueue <- &transaction
	}
}

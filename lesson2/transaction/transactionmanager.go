package transaction

import (
	"log/slog"
	"sync"
	"time"
)

// type WALEntry struct {
// 	entryIndex int64
// 	data       []byte
//  timestamp  int64
// }

type WALEntry []byte

type TransactionManager struct {
	wal              []WALEntry
	snapshot         []byte
	transactionQueue chan []byte
	mutex            sync.Mutex
}

// Создать менеджер транзакций.
func NewTM(transactionQueue chan []byte) *TransactionManager {
	tm := &TransactionManager{
		wal:              make([]WALEntry, 0),
		snapshot:         make([]byte, 0),
		transactionQueue: transactionQueue,
	}

	go func() {
		tm.consumeTransactionsJob()
	}()

	go func() {
		tm.cleanWALJob()
	}()

	return tm
}

func (tm *TransactionManager) cleanWALJob() {
	for range time.Tick(1 * time.Minute) {
		tm.mutex.Lock()

		slog.Info("Cleaning wal, ", "wal", tm.wal, "snapshot", tm.snapshot)
		if len(tm.wal) != 0 {
			tm.snapshot = tm.wal[len(tm.wal)-1]
			tm.wal = tm.wal[:0]
		}

		tm.mutex.Unlock()
	}
}

func (tm *TransactionManager) consumeTransactionsJob() {
	for {
		data := <-tm.transactionQueue

		tm.mutex.Lock()
		tm.wal = append(tm.wal, data)
		tm.mutex.Unlock()
	}
}

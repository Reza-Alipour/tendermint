package raP

import (
	"github.com/tendermint/tendermint/types"
	"sync"
)

var mutex sync.Mutex
var txs = make(map[string][]byte)
var used_txs = make(map[string]bool)

// ReceiveTX done
func ReceiveTX(LastTxHash []byte, Tx []byte) {
	mutex.Lock()
	txs[string(LastTxHash)] = Tx
	mutex.Unlock()
}

// done
func resetTXS() {
	mutex.Lock()
	txs = make(map[string][]byte)
	used_txs = make(map[string]bool)
	mutex.Unlock()
}

// done
func getTX(lastTX types.Tx) ([]byte, bool) {
	mutex.Lock()
	defer mutex.Unlock()

	tx, ok := txs[string(lastTX.Hash())]
	dup_tx := make([]byte, len(tx))
	copy(dup_tx, tx)

	if ok {
		used_txs[string(types.Tx(dup_tx).Hash())] = true
	}
	return dup_tx, ok
}

func getTX_without_lock(lastTX types.Tx) ([]byte, bool) {
	tx, ok := txs[string(lastTX.Hash())]
	dup_tx := make([]byte, len(tx))
	copy(dup_tx, tx)

	if ok {
		used_txs[string(types.Tx(dup_tx).Hash())] = true
	}
	return dup_tx, ok
}

// done
func usedTX(tx types.Tx) bool {
	mutex.Lock()
	defer mutex.Unlock()
	return used_txs[string(tx.Hash())]
}

func usedTX_without_lock(tx types.Tx) bool {
	return used_txs[string(tx.Hash())]
}

func RearrangeTXs(txs types.Txs, maxDataBytes int64) types.Txs {
	infinite := maxDataBytes < 0

	mutex.Lock()

	var new_txs types.Txs
	for _, tx := range txs {
		maxDataBytes -= int64(len(tx))
		if maxDataBytes < 0 && !infinite {
			break
		}
		if !usedTX_without_lock(tx) {
			new_txs = append(new_txs, tx)
		}
		if tx_bytes, ok := getTX_without_lock(tx); ok {
			maxDataBytes -= int64(len(tx))
			if maxDataBytes < 0 && !infinite {
				break
			}
			new_txs = append(new_txs, types.Tx(tx_bytes))
		}
	}
	mutex.Unlock()

	resetTXS()
	return new_txs
}

func GetNextTX(tx types.Tx) types.Tx {
	mutex.Lock()
	defer mutex.Unlock()
	return txs[string(tx.Hash())]
}

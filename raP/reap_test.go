package raP

import (
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/types"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestReceiveTX_1(t *testing.T) {
	lastTx := types.Tx("lastTx 1")
	tx := []byte("transaction 1")

	ReceiveTX(lastTx.Hash(), tx)
	require.Equal(t, tx, txs[string(lastTx.Hash())])
	ac_tx, ok := getTX(lastTx)
	require.Equal(t, true, ok)
	require.Equal(t, tx, ac_tx)
	require.Equal(t, true, usedTX(tx))

	resetTXS()
	require.Equal(t, 0, len(txs))
	require.Equal(t, 0, len(used_txs))
}

func add_some_txs() {
	lastTx := types.Tx("lastTx 1")
	tx := []byte("inserted 1")
	ReceiveTX(lastTx.Hash(), tx)

	lastTx = types.Tx("lastTx 2")
	tx = []byte("inserted 2")
	ReceiveTX(lastTx.Hash(), tx)

	lastTx = types.Tx("lastTx 3")
	tx = []byte("inserted 3")
	ReceiveTX(lastTx.Hash(), tx)

	lastTx = types.Tx("lastTx 4")
	tx = []byte("inserted 4")
	ReceiveTX(lastTx.Hash(), tx)
}

func TestRearrangeTX_1(t *testing.T) {
	add_some_txs()
	tx0 := []byte("transaction 0")
	tx1 := []byte("lastTx 1")
	tx1_1 := []byte("transaction 1_1")
	tx2 := []byte("lastTx 2")
	tx2_1 := []byte("transaction 2_1")
	tx3 := []byte("lastTx 3")
	tx4 := []byte("inserted 1")
	txs_ := types.Txs{tx0, tx1, tx1_1, tx2, tx2_1, tx3, tx4}
	arranged_txs := RearrangeTXs(txs_, -1)
	require.Equal(t, 9, len(arranged_txs))
	require.Equal(t, 0, len(txs))
	require.Equal(t, "transaction 0", string(arranged_txs[0]))
	require.Equal(t, "lastTx 1", string(arranged_txs[1]))
	require.Equal(t, "inserted 1", string(arranged_txs[2]))
	require.Equal(t, "transaction 1_1", string(arranged_txs[3]))
	require.Equal(t, "lastTx 2", string(arranged_txs[4]))
	require.Equal(t, "inserted 2", string(arranged_txs[5]))
	require.Equal(t, "transaction 2_1", string(arranged_txs[6]))
	require.Equal(t, "lastTx 3", string(arranged_txs[7]))
	require.Equal(t, "inserted 3", string(arranged_txs[8]))
}

package coregrpc_test

import (
	"context"
	"github.com/tendermint/tendermint/raP"
	"github.com/tendermint/tendermint/types"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/abci/example/kvstore"
	core_grpc "github.com/tendermint/tendermint/rpc/grpc"
	rpctest "github.com/tendermint/tendermint/rpc/test"
)

func TestMain(m *testing.M) {
	// start a tendermint node in the background to test against
	app := kvstore.NewApplication()
	node := rpctest.StartTendermint(app)

	code := m.Run()

	// and shut down proper at the end
	rpctest.StopTendermint(node)
	os.Exit(code)
}

func TestBroadcastTx(t *testing.T) {
	res, err := rpctest.GetGRPCClient().BroadcastTx(
		context.Background(),
		&core_grpc.RequestBroadcastTx{Tx: []byte("this is a tx")},
	)
	require.NoError(t, err)
	require.EqualValues(t, 0, res.CheckTx.Code)
	require.EqualValues(t, 0, res.DeliverTx.Code)
}

func TestReceiveTX(t *testing.T) {
	tx1 := types.Tx("transaction 0")
	tx2 := types.Tx("inserted transaction")
	_, err := rpctest.GetGRPCClient().RearrangeProposal(
		context.Background(),
		&core_grpc.RequestRearrange{LastTxHash: tx1.Hash(), Tx: tx2},
	)
	require.NoError(t, err)
	tx3 := raP.GetNextTX(tx1)
	require.EqualValues(t, string(tx2), string(tx3))
}

package aggregator

import (
	"fmt"

	"github.com/0xPolygon/cdk/aggregator/prover"
	rpctypes "github.com/0xPolygon/cdk/rpc/types"
	"github.com/ethereum/go-ethereum/common"
)

func hashFromProverPublicInput(b []byte) (common.Hash, error) {
	if len(b) != common.HashLength {
		return common.Hash{}, fmt.Errorf("invalid public input hash length: got %d want %d", len(b), common.HashLength)
	}

	return common.BytesToHash(b), nil
}

func compareFinalProofRootsWithRPC(
	finalProof *prover.FinalProof,
	rpcBatch *rpctypes.RPCBatch,
) (proverSR, rpcSR, proverLER, rpcLER common.Hash, err error) {
	if finalProof == nil {
		return common.Hash{}, common.Hash{}, common.Hash{}, common.Hash{}, fmt.Errorf("final proof is nil")
	}
	if finalProof.Public == nil {
		return common.Hash{}, common.Hash{}, common.Hash{}, common.Hash{}, fmt.Errorf("final proof public inputs are nil")
	}
	if rpcBatch == nil {
		return common.Hash{}, common.Hash{}, common.Hash{}, common.Hash{}, fmt.Errorf("RPC batch is nil")
	}

	proverSR, err = hashFromProverPublicInput(finalProof.Public.NewStateRoot)
	if err != nil {
		return common.Hash{}, common.Hash{}, common.Hash{}, common.Hash{}, fmt.Errorf("prover NewStateRoot: %w", err)
	}
	proverLER, err = hashFromProverPublicInput(finalProof.Public.NewLocalExitRoot)
	if err != nil {
		return common.Hash{}, common.Hash{}, common.Hash{}, common.Hash{}, fmt.Errorf("prover NewLocalExitRoot: %w", err)
	}

	rpcSR = rpcBatch.StateRoot()
	rpcLER = rpcBatch.LocalExitRoot()

	if proverSR != rpcSR || proverLER != rpcLER {
		return proverSR, rpcSR, proverLER, rpcLER, fmt.Errorf(
			"prover/RPC roots mismatch: proverSR=%s rpcSR=%s proverLER=%s rpcLER=%s",
			proverSR, rpcSR, proverLER, rpcLER,
		)
	}

	return proverSR, rpcSR, proverLER, rpcLER, nil
}

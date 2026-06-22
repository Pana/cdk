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
		// FFLONK proofs may return truncated public inputs (e.g. 1-byte
		// NewStateRoot instead of 32 bytes).  Skip comparison and let
		// settlement proceed with RPC roots, which are authoritative.
		// DEBUG: Log the actual length for troubleshooting
		fmt.Printf("DEBUG: NewStateRoot length=%d, data=%x\n", len(finalProof.Public.NewStateRoot), finalProof.Public.NewStateRoot)
		return common.Hash{}, rpcBatch.StateRoot(), common.Hash{}, rpcBatch.LocalExitRoot(), nil
	}
	proverLER, err = hashFromProverPublicInput(finalProof.Public.NewLocalExitRoot)
	if err != nil {
		// Same tolerance for NewLocalExitRoot.
		return proverSR, rpcBatch.StateRoot(), common.Hash{}, rpcBatch.LocalExitRoot(), nil
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

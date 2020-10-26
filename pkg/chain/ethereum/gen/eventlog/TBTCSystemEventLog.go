package eventlog

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	abi "github.com/keep-network/tbtc/pkg/chain/ethereum/gen/abi/system"
	"math/big"
)

// FIXME: This is a temporary structure allowing to access past events
//  emitted by the `TBTCSystem` contract. This structure is here because
//  the generated contract wrappers from `gen/contract` don't support
//  `Filter*` methods yet. When the contract generator will support
//  those methods, the below structure can be removed.
type TBTCSystemEventLog struct {
	contract *abi.TBTCSystem
}

func NewTBTCSystemEventLog(
	contractAddress common.Address,
	backend bind.ContractBackend,
) (*TBTCSystemEventLog, error) {
	contract, err := abi.NewTBTCSystem(contractAddress, backend)
	if err != nil {
		return nil, err
	}

	return &TBTCSystemEventLog{contract}, nil
}

type TBTCSystemRedemptionRequested struct {
	DepositContractAddress common.Address
	Requester              common.Address
	Digest                 [32]byte
	UtxoValue              *big.Int
	RedeemerOutputScript   []byte
	RequestedFee           *big.Int
	Outpoint               []byte
	BlockNumber            uint64
}

func (tsel *TBTCSystemEventLog) PastRedemptionRequestedEvents(
	depositsAddresses []common.Address,
	startBlock uint64,
	endBlock *uint64,
) ([]*TBTCSystemRedemptionRequested, error) {
	iterator, err := tsel.contract.FilterRedemptionRequested(
		&bind.FilterOpts{
			Start: startBlock,
			End:   endBlock,
		},
		depositsAddresses,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}

	events := make([]*TBTCSystemRedemptionRequested, 0)

	for {
		if !iterator.Next() {
			break
		}

		event := iterator.Event
		events = append(events, &TBTCSystemRedemptionRequested{
			DepositContractAddress: event.DepositContractAddress,
			Requester:              event.Requester,
			Digest:                 event.Digest,
			UtxoValue:              event.UtxoValue,
			RedeemerOutputScript:   event.RedeemerOutputScript,
			RequestedFee:           event.RequestedFee,
			Outpoint:               event.Outpoint,
			BlockNumber:            event.Raw.BlockNumber,
		})
	}

	return events, nil
}

package main

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/cockroachdb/pebble"
	abcitypes "github.com/cometbft/cometbft/abci/types"
)

type KVStoreApplication struct {
	db           *pebble.DB
}

var _ abcitypes.Application = (*KVStoreApplication)(nil)

func NewKVStoreApplication(db *pebble.DB) *KVStoreApplication {
	return &KVStoreApplication{db: db}
}
func (app *KVStoreApplication) Info(_ context.Context, info *abcitypes.InfoRequest) (*abcitypes.InfoResponse, error) {
	return &abcitypes.InfoResponse{}, nil
}

func (app *KVStoreApplication) Query(_ context.Context, req *abcitypes.QueryRequest) (*abcitypes.QueryResponse, error) {
	resp := abcitypes.QueryResponse{Key: req.Data}

	val, closer, err := app.db.Get(req.Data)
	if err != nil {
		if !errors.Is(err, pebble.ErrNotFound) {
			log.Panicf("Error reading database, unable to execute query: %e", err)
		}
		resp.Log = "key does not exist"
		return &resp, nil
	}
	resp.Log = "exists"
	resp.Value = val
	closer.Close()
	return &resp, nil
}
func (app *KVStoreApplication) CheckTx(_ context.Context, check *abcitypes.CheckTxRequest) (*abcitypes.CheckTxResponse, error) {
	code := app.isValid(check.Tx)
	return &abcitypes.CheckTxResponse{Code: code}, nil
}

func (app *KVStoreApplication) InitChain(_ context.Context, chain *abcitypes.InitChainRequest) (*abcitypes.InitChainResponse, error) {
	return &abcitypes.InitChainResponse{}, nil
}

func (app *KVStoreApplication) PrepareProposal(_ context.Context, proposal *abcitypes.PrepareProposalRequest) (*abcitypes.PrepareProposalResponse, error) {
	return &abcitypes.PrepareProposalResponse{Txs: proposal.Txs}, nil
}
func (app *KVStoreApplication) ProcessProposal(_ context.Context, proposal *abcitypes.ProcessProposalRequest) (*abcitypes.ProcessProposalResponse, error) {
	return &abcitypes.ProcessProposalResponse{Status: abcitypes.PROCESS_PROPOSAL_STATUS_ACCEPT}, nil
}

func (app *KVStoreApplication) FinalizeBlock(_ context.Context, req *abcitypes.FinalizeBlockRequest) (*abcitypes.FinalizeBlockResponse, error) {
	var txs = make([]*abcitypes.ExecTxResult, len(req.Txs))

	for i, tx := range req.Txs {
		if code := app.isValid(tx); code != 0 {
			log.Printf("Error: invalid transaction index %v", i)
			txs[i] = &abcitypes.ExecTxResult{Code: code}
		} else {
			var transaction Transaction
			if err := transaction.FromBytes(tx); err != nil {
				log.Panicf("Error parsing tx bytes, unable to parse tx: %v", err)
			}
			for _, transfer := range transaction.Transfers {
				src, dst, amount := transfer.Sender, transfer.Dest, transfer.Amount
				log.Printf("Adding key %s with value %s", src, dst)

				srcValue, dstValue := balanceMap[src], balanceMap[dst]

				amountValue, err := strconv.ParseUint(string(amount), 10, 64)
				if err != nil {
					log.Panicf("Error parsing amount, unable to execute tx: %v", err)
				}

				srcValue -= amountValue
				dstValue += amountValue

				if err := app.db.Set([]byte(src), []byte(amount),nil); err != nil{
					log.Panicf("Error writing source key to database, unable to execute tx: %v", err)
				}
				if err := app.db.Set([]byte(dst), []byte(amount),nil); err != nil{
					log.Panicf("Error writing destination key to database, unable to execute tx: %v", err)
				}
				balanceMap[string(src)] = srcValue
				balanceMap[string(dst)] = dstValue
				log.Printf("Successfully added key %s with value %s", src, dst)

				// Add an event for the transaction execution.
				// Multiple events can be emitted for a transaction, but we are adding only one event
				if txs[i] != nil {
					txs[i].Events = append(txs[i].Events, abcitypes.Event{
						Type: "app",
						Attributes: []abcitypes.EventAttribute{
							{Key: "src", Value: string(src), Index: true},
							{Key: "dst", Value: string(dst), Index: true},
							{Key: "amount", Value: string(amount), Index: true},
						},
					})
				} else {
					txs[i] = &abcitypes.ExecTxResult{
						Code: 0,
						Events: []abcitypes.Event{
							{
								Type: "app",
								Attributes: []abcitypes.EventAttribute{
									{Key: "src", Value: string(src), Index: true},
									{Key: "dst", Value: string(dst), Index: true},
									{Key: "amount", Value: string(amount), Index: true},
								},
							},
						},
					}
				}
			}
		}
	}

	return &abcitypes.FinalizeBlockResponse{
		TxResults: txs,
	}, nil
}

func (app KVStoreApplication) Commit(_ context.Context, commit *abcitypes.CommitRequest) (*abcitypes.CommitResponse, error) {
	return &abcitypes.CommitResponse{}, nil
}

func (app *KVStoreApplication) ListSnapshots(_ context.Context, snapshots *abcitypes.ListSnapshotsRequest) (*abcitypes.ListSnapshotsResponse, error) {
	return &abcitypes.ListSnapshotsResponse{}, nil
}

func (app *KVStoreApplication) OfferSnapshot(_ context.Context, snapshot *abcitypes.OfferSnapshotRequest) (*abcitypes.OfferSnapshotResponse, error) {
	return &abcitypes.OfferSnapshotResponse{}, nil
}

func (app *KVStoreApplication) LoadSnapshotChunk(_ context.Context, chunk *abcitypes.LoadSnapshotChunkRequest) (*abcitypes.LoadSnapshotChunkResponse, error) {
	return &abcitypes.LoadSnapshotChunkResponse{}, nil
}

func (app *KVStoreApplication) ApplySnapshotChunk(_ context.Context, chunk *abcitypes.ApplySnapshotChunkRequest) (*abcitypes.ApplySnapshotChunkResponse, error) {
	return &abcitypes.ApplySnapshotChunkResponse{Result: abcitypes.APPLY_SNAPSHOT_CHUNK_RESULT_ACCEPT}, nil
}

func (app KVStoreApplication) ExtendVote(_ context.Context, extend *abcitypes.ExtendVoteRequest) (*abcitypes.ExtendVoteResponse, error) {
	return &abcitypes.ExtendVoteResponse{}, nil
}

func (app *KVStoreApplication) VerifyVoteExtension(_ context.Context, verify *abcitypes.VerifyVoteExtensionRequest) (*abcitypes.VerifyVoteExtensionResponse, error) {
	return &abcitypes.VerifyVoteExtensionResponse{}, nil
}

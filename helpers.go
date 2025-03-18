package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"strconv"
)

type Transfer struct {
	Sender    string `json:"sender"`
	Dest      string `json:"dest"`
	Amount    string `json:"amount"`
	Signature string `json:"signature"`
}

type Transaction struct {
	Transfers []Transfer `json:"transfers"`
}

func (t *Transfer) Challenge() []byte {
	challenge := []byte{}
	challenge = append(challenge, []byte(t.Sender)...)
	challenge = append(challenge, []byte(t.Dest)...)
	challenge = append(challenge, []byte(t.Amount)...)
	return challenge
}

func (t *Transaction) FromBytes(data []byte) error {
	transfersData := bytes.Split(data, []byte(":"))
	for _, transferData := range transfersData {
		parts := bytes.Split(transferData, []byte("="))
		if len(parts) != 4 {
			return errors.New("invalid transaction data")
		}

		transfer := Transfer{
			Sender:    string(parts[0]),
			Dest:      string(parts[1]),
			Amount:    string(parts[2]),
			Signature: string(parts[3]),
		}

		t.Transfers = append(t.Transfers, transfer)
	}
	return nil
}

var keyMap = map[string]string{
	"1": "056d2c2869fb2c1504e80f35f6e85b6b4452c3436e327b4e35e6560ffa95a4c3",
	"2": "4495c8a97d02eecf74e8f319388aee96a902cf7808128577d41d75f600cc2363",
	"3": "61c29c7f938020a0b53358a1726c174f2aa5e8b7685a9ea8b9edb92f7e263ce0",
	"4": "03b5a14fd977cb801e49a8f356b98f8243da388fba9e8653f7d68acbdbebb538",
}

var balanceMap = map[string]uint64{
	"1": 1000000000,
	"2": 1000000000,
	"3": 1000000000,
	"4": 1000000000,
}

func (app *KVStoreApplication) isValid(tx []byte) uint32 {

	var transaction Transaction
	if err := transaction.FromBytes(tx); err != nil {
		return 2
	}

	for _, transfer := range transaction.Transfers {

		if _, ok := keyMap[transfer.Sender]; !ok {
			return 3
		}

		if _, ok := keyMap[transfer.Dest]; !ok {
			return 4
		}
		amount, err := strconv.ParseUint(transfer.Amount, 10, 64)
		if err != nil {
			return 8
		}
		if balanceMap[transfer.Sender] < amount {
			return 5
		}

		pubBytes, err := hex.DecodeString(keyMap[transfer.Sender])
		if err != nil {
			return 6
		}
		pubKey := ed25519.PublicKey(pubBytes)
		signatureBytes, err := hex.DecodeString(transfer.Signature)
		if err != nil {
			return 9
		}
		if !ed25519.Verify(pubKey, transfer.Challenge(), signatureBytes) {
			return 7
		}
	}
	return 0
}

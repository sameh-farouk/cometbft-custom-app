package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"strconv"
)

type Transfer struct {
	Id 	  string `json:"id"`
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
	challenge = append(challenge, []byte(t.Id)...)
	challenge = append(challenge, []byte(t.Sender)...)
	challenge = append(challenge, []byte(t.Dest)...)
	challenge = append(challenge, []byte(t.Amount)...)
	return challenge
}

func (t *Transaction) FromBytes(data []byte) error {
	transfersData := bytes.Split(data, []byte(":"))
	for _, transferData := range transfersData {
		parts := bytes.Split(transferData, []byte("="))
		if len(parts) != 5 {
			return errors.New("invalid transaction data")
		}

		transfer := Transfer{
			Id:        string(parts[0]),
			Sender:    string(parts[1]),
			Dest:      string(parts[2]),
			Amount:    string(parts[3]),
			Signature: string(parts[4]),
		}

		t.Transfers = append(t.Transfers, transfer)
	}
	return nil
}

var keyMap = map[string]string{
	"1": "c8af5ee74756bb934c9c3f93a3ffa4125c93d8a76619a1834f4511334d83d45f",
	"2": "3382d764d3e30ce4c3aab066335a558e8f632d2aaf161e6aa5615c57176cfbca",
	"3": "04c01c7d4f6c784504fce83f97968145e8aa6ca461ec19f3a685466152f17644",
	"4": "d06a22ce4b7a59ceac3a898504901f41e27491ed3cc90e8ee46ac43e9305d61a",
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

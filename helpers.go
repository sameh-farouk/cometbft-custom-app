package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"strconv"
)

type Transaction struct {
	Sender    string `json:"sender"`
	Dest      string `json:"dest"`
	Amount    string `json:"amount"`
	Signature string `json:"signature"`
}

func (t *Transaction) Challenge() []byte {
	challenge := []byte{}
	challenge = append(challenge, []byte(t.Sender)...)
	challenge = append(challenge, []byte(t.Dest)...)
	challenge = append(challenge, []byte(t.Amount)...)
	return challenge
}

func (t *Transaction) FromBytes(data []byte) error {
	parts := bytes.Split(data, []byte("="))
	if len(parts) != 4 {
		return errors.New("invalid transaction data")
	}

	t.Sender = string(parts[0])
	t.Dest = string(parts[1])
	t.Amount = string(parts[2])
	t.Signature = string(parts[3])

	return nil
}

func bytesToUint64(b []byte) uint64 {
	var num uint64
	for i := 0; i < len(b); i++ {
		num |= uint64(b[i]) << (8 * i)
	}
	return num
}

func uint64ToBytes(num uint64) []byte {
	bytes := make([]byte, 8)
	for i := uint(0); i < 8; i++ {
		bytes[i] = byte((num >> (i * 8)) & 0xff)
	}
	return bytes
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
	parts := bytes.Split(tx, []byte("="))
	if len(parts) != 4 {
		return 1
	}

	var t Transaction
	if err := t.FromBytes(tx); err != nil {
		return 2
	}
	if _, ok := keyMap[t.Sender]; !ok {
		return 3
	}

	if _, ok := keyMap[t.Dest]; !ok {
		return 4
	}
	amount, err := strconv.ParseUint(t.Amount, 10, 64)
	if err != nil {
		return 8
	}
	if balanceMap[t.Sender] < amount {
		return 5
	}

	pubBytes, err := hex.DecodeString(keyMap[t.Sender])
	if err != nil {
		return 6
	}
	pubKey := ed25519.PublicKey(pubBytes)
	signatureBytes, err := hex.DecodeString(t.Signature)
	if err != nil {
		return 9
	}
	if !ed25519.Verify(pubKey, t.Challenge(), signatureBytes) {
		return 7
	}

	return 0
}

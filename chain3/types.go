// Copyright (c) 2016, Alan Chen
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package chain3

import (
	"encoding/json"
	"math/big"

	"github.com/caivega/chain3go/common"
)

type JSONBlock struct {
	Difficulty       json.Number `json:"difficulty"`
	ExtraData        string      `json:"extraData"`
	GasLimit         json.Number `json:"gasLimit"`
	GasUsed          json.Number `json:"gasUsed"`
	Hash             string      `json:"hash"`
	LogsBloom        string      `json:"logsBloom"`
	Miner            string      `json:"miner"`
	MixHash          string      `json:"mixHash"`
	Nonce            string      `json:"nonce"`
	Number           json.Number `json:"number"`
	ParentHash       string      `json:"parentHash"`
	ReceiptsRoot     string      `json:"receiptsRoot"`
	Sha3Uncles       string      `json:"sha3Uncles"`
	Size             json.Number `json:"size"`
	StateRoot        string      `json:"stateRoot"`
	Timestamp        json.Number `json:"timestamp"`
	TotalDifficulty  json.Number `json:"totalDifficulty"`
	Transactions     []string    `json:"transactions"`
	TransactionsRoot string      `json:"transactionsRoot"`
	Uncles           []string    `json:"uncles"`
}

func (b *JSONBlock) ToBlock() (block *common.Block) {
	block = &common.Block{}
	block.Difficulty = toBigInt(b.Difficulty)
	block.ExtraData = common.StringToData(b.ExtraData)
	block.GasLimit = toBigInt(b.GasLimit)
	block.GasUsed = toBigInt(b.GasUsed)
	block.Hash = common.StringToHash(b.Hash)
	block.LogsBloom = common.StringToData(b.LogsBloom)
	block.Miner = common.StringToAddress(b.Miner)
	block.MixHash = common.StringToHash(b.MixHash)
	block.Nonce = common.StringToData(b.Nonce)
	block.Number = toBigInt(b.Number)
	block.ParentHash = common.StringToHash(b.ParentHash)
	block.ReceiptsRoot = common.StringToHash(b.ReceiptsRoot)
	block.Sha3Uncles = common.StringToHash(b.Sha3Uncles)
	block.Size = toBigInt(b.Size)
	block.StateRoot = common.StringToHash(b.StateRoot)
	block.Timestamp = toBigInt(b.Timestamp)
	block.TotalDifficulty = toBigInt(b.TotalDifficulty)
	block.Transactions = toHashArray(b.Transactions)
	block.TransactionsRoot = common.StringToHash(b.TransactionsRoot)
	block.Uncles = toHashArray(b.Uncles)
	return block
}

type JSONTransaction struct {
	BlockHash        string      `json:"blockHash"`
	BlockNumber      json.Number `json:"blockNumber"`
	From             string      `json:"from"`
	Gas              json.Number `json:"gas"`
	GasPrice         json.Number `json:"gasprice"`
	Hash             string      `json:"hash"`
	Input            string      `json:"input"`
	Nonce            json.Number `json:"nonce"`
	R                string      `json:"r"`
	S                string      `json:"s"`
	ShardingFlag     string      `json:"shardingFlag"`
	SysCnt           string      `json:"syscnt"`
	To               string      `json:"to"`
	TransactionIndex json.Number `json:"transactionIndex"`
	V                string      `json:"v"`
	Value            json.Number `json:"value"`
}

func (t *JSONTransaction) ToTransaction() (tx *common.Transaction) {
	tx = &common.Transaction{}
	tx.BlockHash = common.StringToHash(t.BlockHash)
	tx.BlockNumber = toBigInt(t.BlockNumber)
	tx.From = common.StringToAddress(t.From)
	tx.Gas = toBigInt(t.Gas)
	tx.GasPrice = toBigInt(t.GasPrice)
	tx.Hash = common.StringToHash(t.Hash)
	tx.Input = common.StringToData(t.Input)
	tx.Nonce = toBigInt(t.Nonce)
	tx.R = common.StringToData(t.R)
	tx.S = common.StringToData(t.S)
	tx.ShardingFlag = common.StringToData(t.ShardingFlag)
	tx.SysCnt = common.StringToData(t.SysCnt)
	tx.To = common.StringToAddress(t.To)
	tx.TransactionIndex = toBigInt(t.TransactionIndex)
	tx.V = common.StringToData(t.V)
	tx.Value = toBigInt(t.Value)
	return tx
}

type JSONTransactionReceipt struct {
	BlockHash         string      `json:"blockHash"`
	BlockNumber       json.Number `json:"blockNumber"`
	ContractAddress   string      `json:"contractAddress"`
	CumulativeGasUsed json.Number `json:"cumulativeGasUsed"`
	From              string      `json:"from"`
	GasUsed           json.Number `json:"gasUsed"`
	Logs              []JSONLog   `json:"logs"`
	LogsBloom         string      `json:"logsBloom"`
	Root              string      `json:"root"`
	To                string      `json:"to"`
	TransactionHash   string      `json:"transactionHash"`
	TransactionIndex  json.Number `json:"transactionIndex"`
}

func (r *JSONTransactionReceipt) ToTransactionReceipt() (receipt *common.TransactionReceipt) {
	receipt = &common.TransactionReceipt{}
	receipt.BlockHash = common.StringToHash(r.BlockHash)
	receipt.BlockNumber = toBigInt(r.BlockNumber)
	receipt.ContractAddress = common.StringToAddress(r.ContractAddress)
	receipt.CumulativeGasUsed = toBigInt(r.CumulativeGasUsed)
	receipt.From = common.StringToAddress(r.From)
	receipt.GasUsed = toBigInt(r.GasUsed)
	receipt.Logs = make([]common.Log, 0)
	for _, l := range r.Logs {
		receipt.Logs = append(receipt.Logs, l.ToLog())
	}
	receipt.LogsBloom = common.StringToData(r.LogsBloom)
	receipt.Root = common.StringToData(r.Root)
	receipt.To = common.StringToAddress(r.To)
	receipt.TransactionHash = common.StringToHash(r.TransactionHash)
	receipt.TransactionIndex = toBigInt(r.TransactionIndex)
	return receipt
}

type JSONLog struct {
	TxData           string      `json:"TxData"`
	Address          string      `json:"address"`
	BlockHash        string      `json:"blockHash"`
	BlockNumber      json.Number `json:"blockNumber"`
	LogIndex         json.Number `json:"logIndex"`
	Removed          bool        `json:"removed"`
	Topics           []string    `json:"topics"`
	TransactionHash  string      `json:"transactionHash"`
	TransactionIndex json.Number `json:"transactionIndex"`
}

func (l JSONLog) ToLog() (log common.Log) {
	log = common.Log{}
	log.TxData = common.StringToData(l.TxData)
	log.Address = common.StringToAddress(l.Address)
	log.BlockHash = common.StringToHash(l.BlockHash)
	log.BlockNumber = toBigInt(l.BlockNumber)
	log.LogIndex = toBigInt(l.LogIndex)
	log.Removed = l.Removed
	log.Topics = toDataArray(l.Topics)
	log.TransactionHash = common.StringToHash(l.TransactionHash)
	log.TransactionIndex = toBigInt(l.TransactionIndex)
	return log
}

func toBigInt(data json.Number) *big.Int {
	f := new(big.Int)
	f.SetString(string(data), 0)
	return f
}

func toHashArray(list []string) []common.Hash {
	c := len(list)
	hashs := make([]common.Hash, c)
	for i := 0; i < c; i++ {
		hashs[i] = common.StringToHash(list[i])
	}
	return hashs
}

func toDataArray(list []string) []common.Data {
	c := len(list)
	hashs := make([]common.Data, c)
	for i := 0; i < c; i++ {
		hashs[i] = common.StringToData(list[i])
	}
	return hashs
}

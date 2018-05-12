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

package common

import (
	"encoding/json"
	"math/big"
)

const (
	hashLength    = 32
	addressLength = 20
)

// Hash ...
type Hash [hashLength]byte

func NewHash(bytes []byte) (result Hash) {
	copy(result[:], bytes)
	return result
}

func (hash *Hash) String() string {
	return BytesToHex(hash[:])
}

// Data ...
type Data []byte

func NewData(bytes []byte) (data Data) {
	copy(data[:], bytes)
	return data
}

func (data *Data) String() string {
	return BytesToHex((*data)[:])
}

// Address ...
type Address [addressLength]byte

func NewAddress(bytes []byte) (result Address) {
	copy(result[:], bytes)
	return result
}

func (addr *Address) String() string {
	return BytesToHex(addr[:])
}

// SyncStatus ...
type SyncStatus struct {
	Result        bool
	StartingBlock *big.Int
	CurrentBlock  *big.Int
	HighestBlock  *big.Int
}

// TransactionRequest ...
type TransactionRequest struct {
	From     Address  `json:"from"`
	To       Address  `json:"to"`
	Gas      *big.Int `json:"gas"`
	GasPrice *big.Int `json:"gasprice"`
	Value    *big.Int `json:"value"`
	Data     []byte   `json:"data"`
}

func (tx *TransactionRequest) String() string {
	jsonBytes, _ := json.Marshal(tx)
	return string(jsonBytes)
}

// Transaction ...
type Transaction struct {
	BlockHash        Hash     `json:"blockHash"`
	BlockNumber      *big.Int `json:"blockNumber"`
	From             Address  `json:"from"`
	Gas              *big.Int `json:"gas"`
	GasPrice         *big.Int `json:"gasprice"`
	Hash             Hash     `json:"hash"`
	Input            Data     `json:"input"`
	Nonce            *big.Int `json:"nonce"`
	R                Data     `json:"r"`
	S                Data     `json:"s"`
	ShardingFlag     Data     `json:"shardingFlag"`
	SysCnt           Data     `json:"syscnt"`
	To               Address  `json:"to"`
	TransactionIndex *big.Int `json:"transactionIndex"`
	V                Data     `json:"v"`
	Value            *big.Int `json:"value"`
}

func (tx *Transaction) String() string {
	jsonBytes, _ := json.Marshal(tx)
	return string(jsonBytes)
}

// Log ...
type Log struct {
	TxData           Data     `json:"TxData"`
	Address          Address  `json:"address"`
	BlockHash        Hash     `json:"blockHash"`
	BlockNumber      *big.Int `json:"blockNumber"`
	LogIndex         *big.Int `json:"logIndex"`
	Removed          bool     `json:"removed"`
	Topics           []Data   `json:"topics"`
	TransactionHash  Hash     `json:"transactionHash"`
	TransactionIndex *big.Int `json:"transactionIndex"`
}

// TransactionReceipt ...
type TransactionReceipt struct {
	BlockHash         Hash     `json:"blockHash"`
	BlockNumber       *big.Int `json:"blockNumber"`
	ContractAddress   Address  `json:"contractAddress"`
	CumulativeGasUsed *big.Int `json:"cumulativeGasUsed"`
	From              Address  `json:"from"`
	GasUsed           *big.Int `json:"gasUsed"`
	Logs              []Log    `json:"logs"`
	LogsBloom         Data     `json:"logsBloom"`
	Root              Data     `json:"root"`
	To                Address  `json:"to"`
	TransactionHash   Hash     `json:"transactionHash"`
	TransactionIndex  *big.Int `json:"transactionIndex"`
}

func (tx *TransactionReceipt) String() string {
	jsonBytes, _ := json.Marshal(tx)
	return string(jsonBytes)
}

// Block ...
type Block struct {
	Difficulty       *big.Int `json:"difficulty"`
	ExtraData        Data     `json:"extraData"`
	GasLimit         *big.Int `json:"gasLimit"`
	GasUsed          *big.Int `json:"gasUsed"`
	Hash             Hash     `json:"hash"`
	LogsBloom        Data     `json:"logsBloom"`
	Miner            Address  `json:"miner"`
	MixHash          Hash     `json:"mixHash"`
	Nonce            Data     `json:"nonce"`
	Number           *big.Int `json:"number"`
	ParentHash       Hash     `json:"parentHash"`
	ReceiptsRoot     Hash     `json:"receiptsRoot"`
	Sha3Uncles       Hash     `json:"sha3Uncles"`
	Size             *big.Int `json:"size"`
	StateRoot        Hash     `json:"stateRoot"`
	Timestamp        *big.Int `json:"timestamp"`
	TotalDifficulty  *big.Int `json:"totalDifficulty"`
	Transactions     []Hash   `json:"transactions"`
	TransactionsRoot Hash     `json:"transactionsRoot"`
	Uncles           []Hash   `json:"uncles"`
}

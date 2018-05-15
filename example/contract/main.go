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

package main

import (
	"flag"
	"fmt"

	Chain3 "github.com/caivega/chain3go/chain3"
	"github.com/caivega/chain3go/common"
	"github.com/caivega/chain3go/provider"
	"github.com/caivega/chain3go/rpc"
)

var hostname = flag.String("hostname", "localhost", "The ethereum client RPC host")
var port = flag.String("port", "8545", "The ethereum client RPC port")
var contract = flag.String("contract", "0x8b14c0f1de8159204f841850606bf6ac36ed89b3", "default filter address")
var verbose = flag.Bool("verbose", false, "Print verbose messages")

func main() {
	flag.Parse()

	if *verbose {
		fmt.Printf("Connect to %s:%s\n", *hostname, *port)
	}

	provider := provider.NewHTTPProvider(*hostname+":"+*port, rpc.GetDefaultMethod())
	chain3 := Chain3.NewChain3(provider)
	mc := chain3.Mc

	filter, err := mc.NewBlockFilter()
	if err != nil {
		fmt.Printf("Failed to create filter, %v\n", err)
		return
	}
	defer func() {
		if _, err := mc.UninstallFilter(filter); err != nil {
			fmt.Printf("UninstallFilter failed: %v\n", err)
		}

	}()

	fmt.Printf("Filter ID: %v\n", filter.ID())

	if filterCh := filter.Watch(); filterCh != nil {
		for {
			log, err := filterCh.Next()
			if err == nil {
				blockHash := log.(string)
				block, err := mc.GetBlockByHash(common.StringToHash(blockHash), false)
				if err != nil {
					fmt.Println("error", err)
				}
				for i := 0; i < len(block.Transactions); i++ {
					tx, err := mc.GetTransactionByHash(block.Transactions[i])
					if err != nil {
						fmt.Println("error", err)
					}
					if tx.From.String() == *contract || tx.To.String() == *contract {
						fmt.Println("tx", i, tx.Hash.String(), tx.Input.String())
					}
				}
			} else {
				fmt.Println("error", err)
				return
			}
		}
	}
}

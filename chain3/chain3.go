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
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/caivega/chain3go/common"
	"github.com/caivega/chain3go/provider"
	"github.com/tonnerre/golang-go.crypto/sha3"
)

var (
	big0    = big.NewInt(0)
	rat0    = big.NewRat(0, 1)
	unitMap = map[string]string{
		"Gsha":    "1000000000",
		"Ksha":    "1000",
		"Msha":    "1000000",
		"femtomc": "1000",
		"gmc":     "1000000000000000000000000000",
		"grand":   "1000000000000000000000",
		"gsha":    "1000000000",
		"kmc":     "1000000000000000000000",
		"ksha":    "1000",
		"mc":      "1000000000000000000",
		"micro":   "1000000000000",
		"micromc": "1000000000000",
		"milli":   "1000000000000000",
		"millimc": "1000000000000000",
		"mmc":     "1000000000000000000000000",
		"msha":    "1000000",
		"nano":    "1000000000",
		"nanomc":  "1000000000",
		"nomc":    "0",
		"picomc":  "1000000",
		"sand":    "1000000000000",
		"sha":     "1",
		"tmc":     "1000000000000000000000000000000",
		"xiao":    "1000000000",
	}
)

// Chain3 Standard interface
// See https://github.com/ethereum/wiki/wiki/JavaScript-API#web3js-api-reference
type Chain3 struct {
	provider       provider.Provider
	requestManager *RequestManager
	Mc             Mc
	Net            Net
}

// NewChain3 creates a new chain3 object.
func NewChain3(provider provider.Provider) *Chain3 {
	requestManager := NewRequestManager(provider)
	return &Chain3{
		provider:       provider,
		requestManager: requestManager,
		Mc:             newMoacAPI(requestManager),
		Net:            newNetAPI(requestManager)}
}

// IsConnected checks if a connection to a node exists.
func (chain3 *Chain3) IsConnected() bool {
	return true
}

// SetProvider sets provider.
func (chain3 *Chain3) SetProvider(provider provider.Provider) {
	chain3.provider = provider
}

// CurrentProvider returns the current provider.
func (chain3 *Chain3) CurrentProvider() provider.Provider {
	return chain3.provider
}

func (chain3 *Chain3) CurrentRequestManager() *RequestManager {
	return chain3.requestManager
}

// Reset state of chain3. Resets everything except manager. Uninstalls all
// filters. Stops polling. If keepSyncing is true, it will uninstall all
// filters, but will keep the chain3.mc.IsSyncing() polls.
func (chain3 *Chain3) Reset(keepSyncing bool) {

}

// Sha3 returns Keccak-256 (not the standardized SHA3-256) of the given data.
func (chain3 *Chain3) Sha3(data string, options interface{}) string {
	opt := struct {
		Encoding string `json:"encoding"`
	}{
		"default",
	}

checkEncoding:
	switch options.(type) {
	case string:
		if err := json.Unmarshal([]byte(options.(string)), &opt); err != nil {
			return common.BytesToHex(chain3.sha3Hash([]byte(data)))
		}
		break checkEncoding
	default:
		var err error
		var optBytes []byte
		if optBytes, err = json.Marshal(options); err != nil {
			return common.BytesToHex(chain3.sha3Hash([]byte(data)))
		}

		if err = json.Unmarshal(optBytes, &opt); err != nil {
			return common.BytesToHex(chain3.sha3Hash([]byte(data)))
		}
		break checkEncoding
	}

	if opt.Encoding == "hex" {
		return common.BytesToHex(chain3.sha3Hash(common.HexToBytes(data)))
	}
	return common.BytesToHex(chain3.sha3Hash([]byte(data)))
}

// ToHex converts any value into HEX.
func (chain3 *Chain3) ToHex(value interface{}) string {
	switch value.(type) {
	case bool:
		v := value.(bool)
		if v {
			return "0x1"
		}
		return "0x0"
	case string:
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return chain3.FromDecimal(value)
		}
		unquoted, err := strconv.Unquote(string(jsonBytes))
		if err != nil {
			return common.BytesToHex(jsonBytes)
		}
		return common.BytesToHex([]byte(unquoted))
	case *big.Int:
		return chain3.FromDecimal(value)
	default:
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return chain3.FromDecimal(value)
		}
		return common.BytesToHex(jsonBytes)
	}
}

// ToASCII converts a HEX string into a ASCII string.
func (chain3 *Chain3) ToASCII(hexString string) string {
	return string(bytes.Trim(common.HexToBytes(hexString), "\x00"))
}

// FromASCII converts any ASCII string to a HEX string.
func (chain3 *Chain3) FromASCII(textString string, padding int) string {
	hex := ""
	for _, runeValue := range textString {
		hex += fmt.Sprintf("%x", runeValue)
	}

	l := len(hex)
	for i := 0; i < padding*2-l; i++ {
		hex += "0"
	}
	return "0x" + hex
}

// ToDecimal converts value to it"s decimal representation in string.
func (chain3 *Chain3) ToDecimal(value interface{}) string {
	n := chain3.ToBigNumber(value)
	if n.IsInt() {
		return n.Num().String()
	}
	return n.String()
}

// FromDecimal converts value to it"s hex representation.
func (chain3 *Chain3) FromDecimal(value interface{}) string {
	number := chain3.ToBigNumber(value)
	if number.IsInt() {
		result := number.Num().Text(16)

		if number.Cmp(rat0) < 0 {
			return "-0x" + result[1:]
		}
		return "0x" + result
	}

	v, _ := number.Float64()
	return fmt.Sprintf("%x", math.Float64bits(v))
}

// FromWei takes a number of wei and converts it to any other ether unit.
//
// Possible units are:
//   SI Short   SI Full        Effigy       Other
// - kwei       femtoether     babbage
// - mwei       picoether      lovelace
// - gwei       nanoether      shannon      nano
// - --         microether     szabo        micro
// - --         microether     szabo        micro
// - --         milliether     finney       milli
// - ether      --             --
// - kether                    --           grand
// - mether
// - gether
// - tether
func (chain3 *Chain3) FromSha(number string, unit string) string {
	num := chain3.ToBigNumber(number)
	returnValue := num.Quo(num, chain3.getValueOfUnit(unit))
	return returnValue.Num().String()
}

// ToWei takes a number of a unit and converts it to wei.
//
// Possible units are:
//   SI Short   SI Full        Effigy       Other
// - kwei       femtoether     babbage
// - mwei       picoether      lovelace
// - gwei       nanoether      shannon      nano
// - --         microether     szabo        micro
// - --         microether     szabo        micro
// - --         milliether     finney       milli
// - ether      --             --
// - kether                    --           grand
// - mether
// - gether
// - tether
func (chain3 *Chain3) ToSha(number interface{}, unit string) string {
	num := chain3.ToBigNumber(number)
	returnValue := num.Mul(num, chain3.getValueOfUnit(unit))
	return returnValue.Num().String()
}

// ToBigNumber takes an input and transforms it into an *big.Rat.
func (chain3 *Chain3) ToBigNumber(value interface{}) (result *big.Rat) {
	switch value.(type) {
	case *big.Rat:
		v := value.(*big.Rat)
		return v
	case *big.Int:
		v := value.(*big.Int)
		result = new(big.Rat)
		result.SetInt(v)
		return result
	case string:
		v := value.(string)
		i := new(big.Int)
		result = new(big.Rat)

		if strings.Index(v, "0x") == 0 || strings.Index(v, "-0x") == 0 {
			i.SetString(strings.Replace(v, "0x", "", -1), 16)
		} else {
			i.SetString(v, 10)
		}
		result.SetInt(i)
		return result
	}
	return result
}

// IsAddress checks if the given string is an address.
func (chain3 *Chain3) IsAddress(address string) bool {
	smallCapsMatcher := regexp.MustCompile("^(0x)?[0-9a-f]{40}$")
	smallCapsMatched := smallCapsMatcher.MatchString(address)
	allCapsMatcher := regexp.MustCompile("^(0x)?[0-9A-F]{40}$")
	allCapsMatched := allCapsMatcher.MatchString(address)
	if smallCapsMatched || allCapsMatched {
		return true
	}
	return chain3.isChecksumAddress(address)
}

func (chain3 *Chain3) isChecksumAddress(address string) bool {
	addr := strings.Replace(address, "0x", "", -1)
	addressHash := chain3.Sha3(strings.ToLower(addr), "")

	for i := 0; i < 40; i++ {
		d, err := strconv.ParseInt(string(addressHash[i]), 16, 32)
		if err != nil {
			return false
		}

		if d > 7 && strings.ToUpper(string(address[i])) == string(address[i]) ||
			d <= 7 && strings.ToLower(string(address[i])) == string(address[i]) {
			return false
		}
	}
	return true
}

func (chain3 *Chain3) sha3Hash(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}

func (chain3 *Chain3) getValueOfUnit(unit string) *big.Rat {
	u := strings.TrimSpace(unit)
	if u != "" {
		u = strings.ToLower(u)
	} else {
		u = "ether"
	}

	if unitValue, ok := unitMap[u]; ok {
		value := new(big.Int)
		value.SetString(unitValue, 10)
		returnValue := new(big.Rat)
		returnValue.SetInt(value)
		return returnValue
	}

	keys := make([]string, 0, len(unitMap))
	for k := range unitMap {
		keys = append(keys, k)
	}
	panic(fmt.Sprintf("This unit doesn\"t exists, please use the one of the following units, %v", keys))
}

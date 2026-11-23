/*
Copyright © 2021 Billy G. Allie <bill.allie@defiant.mug.org>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"math/bits"
	"os"
	"strings"

	"github.com/bgallie/ikmachine/v2"
	"github.com/bgallie/jc1/v2"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	// rRead is a function variable that can be set to use either the crypto/rand or ikmachine read function.
	rRead func([]byte) (int, error)

	// rPerm is a function variable that can be set to use either the crypto/rand or the ikmachine perm function.
	rPerm func(int) []int

	// rInt is a function variable that can be set to use either the crypto/rand or ikmachine Intn function.
	rInt func(int64) int64
)

func cryptographicRandom() {
	// Define rRead and rInt to use the crypto/rand based functions.
	rRead = rand.Read
	rPerm = perm
	rInt = cInt
}

func cInt(n int64) int64 {
	max := big.NewInt(n)
	j, err := rand.Int(rand.Reader, max)
	cobra.CheckErr(err)
	return j.Int64()
}

func perm(n int) []int {
	if n < 0 {
		panic(fmt.Sprintf("Perm called with a negative argument [%d]", n))
	}
	res := make([]int, n)
	for i := 1; i < n; i++ {
		j := rInt(int64(i + 1))
		res[i] = res[j]
		res[j] = i
	}
	return res
}

func ikmachineRandom(args []string) error {
	// Get the secret key from either the environment variable or the command line arguments.
	secret := getKeyFromEnvOrArgs(args)

	// Initialize the ikmachine with the secret key and the named proforma file.
	ikEngine := new(ikmachine.IkMachine).InitializeProformaEngine().ApplyKey('E', []byte(secret))
	// Get the counter key by setting the index to zero (0), encrypting an
	// cipherblock of all zeros, encoding the results in base64.
	savedIndex := ikEngine.Index() // save the current index.
	ikEngine.Index(ikmachine.BigZero)
	ikEngine.Input <- make(ikmachine.CipherBlock, ikmachine.CipherBlockBytes)
	counterKey = base64.RawStdEncoding.EncodeToString(<-ikEngine.Output)
	ikEngine.Index(savedIndex) // restore the saved index.
	// Get the random functions
	ikRandom := new(ikmachine.Rand).New(ikEngine)
	shutdownEngine = func() {
		ikRandom.StopRand()
		ikRandom = nil
		cobra.CheckErr(writeCounterToConfig(counterKey, ikEngine.Index()))
		ikEngine.StopIkMachine()
		ikEngine = nil
	}
	var iCnt *big.Int = new(big.Int).Set(ikmachine.BigZero)
	// Get the starting block count.  cnt can be a number or a fraction such
	// as "1/2", "2/3", or "3/4".  If it is a fraction, then the starting block
	// count is calculated by multiplying the maximal states of the ikengine
	// by the fraction.
	if bCnt != "" {
		var good bool
		flds := strings.Split(bCnt, "/")
		if len(flds) == 1 {
			iCnt, good = new(big.Int).SetString(bCnt, 10)
			if !good {
				return fmt.Errorf("failed converting the count to a big.Int: [%s]", bCnt)
			}
		} else if len(flds) == 2 {
			m := new(big.Int).Set(ikEngine.MaximalStates())
			a, good := new(big.Int).SetString(flds[0], 10)
			if !good {
				return fmt.Errorf("failed converting the numerator to a big.Int: [%s]", flds[0])
			}
			b, good := new(big.Int).SetString(flds[1], 10)
			if !good {
				return fmt.Errorf("failed converting the denominator to a big.Int: [%s]", flds[1])
			}
			iCnt = m.Div(m.Mul(m, a), b)
		} else {
			return fmt.Errorf("incorrect initial count: [%s]", bCnt)
		}
	} else {
		// Read the saved counter from the config file.
		var err error = nil
		iCnt, err = readCounterFromConfig(counterKey, iCnt)
		cobra.CheckErr(err)
	}

	// Now we can set the index of the cipher machine.
	ikEngine.Index(iCnt)

	// Set the random functions to use the ikmachine random functions.
	rRead = ikRandom.Read
	rPerm = ikRandom.Perm
	rInt = ikRandom.Int63n

	return nil
}

func jc1Random(args []string) {
	// Get the secret key from either the environment variable or the command line arguments.
	secret := getKeyFromEnvOrArgs(args)

	// Initialize the JC1 Cipher with the secret key.
	cipher := new(jc1.Cipher).New([]byte(secret))
	// Get the random functions
	rRead = cipher.Read
	rPerm = perm
	rInt = int63n
}

func uberjc1Random(args []string) {
	// Get the secret key from either the environment variable or the command line arguments.
	secret := getKeyFromEnvOrArgs(args)

	// Initialize the JC1 Cipher with the secret key.
	cipher := new(jc1.UberJc1).New([]byte(secret))
	// Get the random functions
	rRead = cipher.Read
	rPerm = perm
	rInt = int63n
}

// int63n returns, as an int64, a non-negative pseudo-random number in the
// half-open interval [0,n). It panics if n <= 0.
func int63n(max int64) int64 {
	if max <= 0 {
		panic("argument to Int63n is <= 0")
	}
	return int64(uint64n(uint64(max)))
}

// uint64n returns, as an uint64, a unsigned 64 bit pseudo-random number in
// the half-open interval [0, max).
func uint64n(max uint64) uint64 {
	if max <= 0 {
		panic("argument to Int63n is <= 0")
	}
	n := max - 1
	// bitLen is the maximum bit length needed to encode a value < max.
	bitLen := bits.Len64(n)
	if bitLen == 0 {
		// the only valid result is 0
		return n
	}
	// k is the maximum byte length needed to encode a value < max.
	k := (bitLen + 7) / 8
	// b is the number of bits in the most significant byte of max-1.
	b := uint64(bitLen % 8)
	if b == 0 {
		b = 8
	}
	bytes := make([]byte, k)
	for {
		_, _ = rRead(bytes[:0])
		// Clear bits in the first byte to increase the probability
		// that the candidate is < max.
		bytes[0] &= uint8((1 << b) - 1)
		// Change the data in the byte slice into an integer ('n')
		n = 0
		for _, val := range bytes {
			n = (n << 8) | uint64(val)
		}
		if n < max {
			return n
		}
	}
}

func getKeyFromEnvOrArgs(args []string) string {
	// Obtain the passphrase used to encrypt the file from either:
	// 1. User input from the terminal (most secure)
	// 2. The 'TNT2_SECRET' environment variable (less secure)
	// 3. Arguments from the entered command line (least secure - not recommended)
	var secret string
	if len(args) == 0 {
		if grViper.IsSet("GPF_SECRET") {
			secret = grViper.GetString("GPF_SECRET")
		} else {
			if term.IsTerminal(int(os.Stdin.Fd())) {
				fmt.Fprintf(os.Stderr, "Enter the passphrase: ")
				byteSecret, err := term.ReadPassword(int(os.Stdin.Fd()))
				cobra.CheckErr(err)
				fmt.Fprintln(os.Stderr, "")
				secret = string(byteSecret)
			}
		}
	} else {
		secret = strings.Join(args, " ")
	}

	if len(secret) == 0 {
		cobra.CheckErr("You must supply a password.")
	}
	return secret
}

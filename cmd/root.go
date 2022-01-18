/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"os/user"
	"regexp"
	"strings"

	"github.com/bgallie/tntengine"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/term"

	"github.com/spf13/viper"
)

var (
	cfgFile    string
	multiplier = map[string]*big.Int{
		"":   new(big.Int).SetInt64(1),
		"c":  new(big.Int).SetInt64(1),
		"w":  new(big.Int).SetInt64(2),
		"b":  new(big.Int).SetInt64(512),
		"K":  new(big.Int).SetInt64(1000),
		"KB": new(big.Int).SetInt64(1024),
		"M":  new(big.Int).SetInt64(1000 * 1000),
		"MB": new(big.Int).SetInt64(1024 * 1024),
		"G":  new(big.Int).SetInt64(1000 * 1000 * 1000),
		"GB": new(big.Int).SetInt64(1024 * 1024 * 1024),
		"T":  new(big.Int).SetInt64(1000 * 1000 * 1000 * 1000),
		"TB": new(big.Int).SetInt64(1024 * 1024 * 1024 * 1024),
		"P":  new(big.Int).Mul(new(big.Int).SetInt64(1000*1000*1000*1000), new(big.Int).SetInt64(1000)),
		"PB": new(big.Int).Mul(new(big.Int).SetInt64(1024*1024*1024*1024), new(big.Int).SetInt64(1024)),
		"E":  new(big.Int).Mul(new(big.Int).SetInt64(1000*1000*1000*1000), new(big.Int).SetInt64(1000*1000)),
		"EB": new(big.Int).Mul(new(big.Int).SetInt64(1024*1024*1024*1024), new(big.Int).SetInt64(1024*1024)),
		"Z":  new(big.Int).Mul(new(big.Int).SetInt64(1000*1000*1000*1000), new(big.Int).SetInt64(1000*1000*1000)),
		"ZB": new(big.Int).Mul(new(big.Int).SetInt64(1024*1024*1024*1024), new(big.Int).SetInt64(1024*1024*1024)),
		"Y":  new(big.Int).Mul(new(big.Int).SetInt64(1000*1000*1000*1000), new(big.Int).SetInt64(1000*1000*1000*1000)),
		"YB": new(big.Int).Mul(new(big.Int).SetInt64(1024*1024*1024*1024), new(big.Int).SetInt64(1024*1024*1024*1024)),
	}
	proFormaFileName string
	tntMachine       tntengine.TntEngine
	bCnt             string
	iCnt             *big.Int
	sCnt             string
	sBlock           string
	count            *big.Int
	blockSize        *big.Int
	cMap             map[string]*big.Int
	mKey             string
	outputFileName   string
	cntrFileName     string
	GitCommit        string = "not set"
	GitBranch        string = "not set"
	GitState         string = "not set"
	GitSummary       string = "not set"
	BuildDate        string = "not set"
	Version          string = "dev"
)

const (
	genRandCountFile = ".genRand"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "genRandom",
	Short: "Generate (psudo)random data using the tnt2engine random number generator",
	Long: `Generate (psudo)random data using the tnt2engine random number generator as either
	a stream of random bytes, a stream of ASCII '0' and '1' characters, a stream of hexidecimal encoded bytes,
	or a series of X,Y coordinatges within a rectangle of a given size.`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tnt2.yaml)")
	rootCmd.PersistentFlags().StringVarP(&bCnt, "count", "", "", `starting block count.
The count can be a number or a fraction such as "1/2", "2/3", or "3/4".  If it is a fraction, then the starting 
block count is calculated by multiplying the maximal blocks generated by the tntEngine by the fraction.
Supplying a count will overide the stored count in the .genRand file, allowing for a repeatable stream of
psuedo random data by giving the same secret key and starting block number.`)
	rootCmd.PersistentFlags().StringVarP(&proFormaFileName, "proformafile", "f", "", "the file name containing the proforma machine to use instead of the builtin proforma machine.")
	rootCmd.PersistentFlags().StringVarP(&outputFileName, "outputFile", "o", "-", "Name of the file containing the generated (psudo)random data.")
	// rootCmd.MarkPersistentFlagRequired("outputFile")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".genRandom" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".genRandom")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Get the counter file name based on the current user.
	u, err := user.Current()
	cobra.CheckErr(err)
	cntrFileName = fmt.Sprintf("%s%c%s", u.HomeDir, os.PathSeparator, genRandCountFile)
}

// ParseNumber parses a numeric string with an optional multiplier into a number.
func ParseNumber(num string) *big.Int {
	val := new(big.Int)
	re := regexp.MustCompile(`([[:digit:]]+)([[:alpha:]]*)`)
	tokens := re.FindSubmatch([]byte(num))
	if tokens == nil {
		fmt.Fprintf(os.Stderr, "Bad numeric string format: [%s]\n", num)
		os.Exit(1)
	}
	_, ok := val.SetString(string(tokens[1]), 0)
	if !ok {
		fmt.Fprintf(os.Stderr, "Failed to convert string to a big.Int: %s\n", tokens[1])
		os.Exit(1)
	}
	mul := multiplier[string(tokens[2])]
	if mul == nil {
		fmt.Fprintf(os.Stderr, "Failed to convert string to a multiplier: %s\n", tokens[2])
		os.Exit(1)
	}
	return val.Mul(val, mul)
}

func initEngine(args []string) {
	// Obtain the passphrase used to encrypt the file from either:
	// 1. User input from the terminal (most secure)
	// 2. The 'TNT2_SECRET' environment variable (less secure)
	// 3. Arguments from the entered command line (least secure - not recommended)
	var secret string
	if len(args) == 0 {
		if viper.IsSet("TNT2_SECRET") {
			secret = viper.GetString("TNT2_SECRET")
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
		// } else {
		// 	fmt.Printf("Secret: [%s]\n", secret)
	}

	// Initialize the tntengine with the secret key and the named proforma file.
	tntMachine.Init([]byte(secret), proFormaFileName)
	// Set the engine type and build the cipher machine.
	tntMachine.SetEngineType("E")
	tntMachine.BuildCipherMachine()
	// Get the starting block count.  cnt can be a number or a fraction such
	// as "1/2", "2/3", or "3/4".  If it is a fraction, then the starting block
	// count is calculated by multiplying the maximal states of the tntEngine
	// by the fraction.
	if bCnt != "" {
		var good bool
		flds := strings.Split(bCnt, "/")
		if len(flds) == 1 {
			iCnt, good = new(big.Int).SetString(bCnt, 10)
			if !good {
				cobra.CheckErr(fmt.Sprintf("Failed converting the count to a big.Int: [%s]\n", bCnt))
			}
		} else if len(flds) == 2 {
			m := new(big.Int).Set(tntMachine.MaximalStates())
			a, good := new(big.Int).SetString(flds[0], 10)
			if !good {
				cobra.CheckErr(fmt.Sprintf("Failed converting the numerator to a big.Int: [%s]\n", flds[0]))
			}
			b, good := new(big.Int).SetString(flds[1], 10)
			if !good {
				cobra.CheckErr(fmt.Sprintf("Failed converting the denominator to a big.Int: [%s]\n", flds[1]))
			}
			iCnt = m.Div(m.Mul(m, a), b)
		} else {
			cobra.CheckErr(fmt.Sprintf("Incorrect initial count: [%s]\n", bCnt))
		}
	} else {
		iCnt = new(big.Int).Set(tntengine.BigZero)
	}

	// Read in the map of counts from the file which holds the counts and get
	// the count to use to encode the file.
	cMap = make(map[string]*big.Int)
	cMap = readCounterFile(cMap)
	mKey = tntMachine.CounterKey()
	if cMap[mKey] == nil {
		cMap[mKey] = iCnt
	} else {
		if bCnt != "" {
			fmt.Fprintf(os.Stderr, "Overriding the value from the %s file.\n", genRandCountFile)
		} else {
			iCnt = cMap[mKey]
		}
	}

	// Now we can set the index of the ciper machine.
	tntMachine.SetIndex(iCnt)
}

// shutdownTnTMachine will write the current block count to the counter file
// and shut down the TntMachine (by exiting the go functions in the TntMachine).
func shutdownTntMachine() {
	cMap[mKey] = tntMachine.Index()
	checkError(writeCounterFile(cMap))

	var blk tntengine.CypherBlock
	tntMachine.Left() <- blk
	<-tntMachine.Right()
}

/*
	getOutputFiles will return the output file to use while generating the
	random data.  If an output file names was given, then that file will be
	opened.  Otherwise stdout is used.
*/
func getOutputFile() *os.File {
	var err error
	var fout *os.File

	if len(outputFileName) > 0 {
		if outputFileName == "-" {
			fout = os.Stdout
		} else {
			fout, err = os.Create(outputFileName)
			cobra.CheckErr(err)
		}
	} else {
		fout = os.Stdout
	}

	return fout
}

// generatRandomStream write the (psudo)random data generate by the tntengine to a
// io.Pipe which can be read from the returned io.PipeReader.  it will generate and
// write (blockSize X blockCount) (psudo)random bytes.
func generateRandomStream() *io.PipeReader {
	count = ParseNumber(sCnt)
	blockSize = ParseNumber(sBlock)
	input, output := io.Pipe()
	if !blockSize.IsInt64() || blockSize.Int64() > multiplier["MB"].Int64() {
		fmt.Fprintln(os.Stderr, "Block size must be less than or equal to 1,000,000")
		os.Exit(1)
	}
	blkSize := blockSize.Int64()
	block := make([]byte, blkSize)
	q, r := new(big.Int).QuoRem(count, big.NewInt(math.MaxInt64), new(big.Int))
	first, last := int64(0), q.Int64()
	random := tntengine.NewRand(&tntMachine)

	go func() {
		defer output.Close()
		blocksWritten := new(big.Int).Set(tntengine.BigZero)
		for i := first; i < last; i++ {
			for j := first; j < math.MaxInt64; j++ {
				_, err := random.Read(block)
				cobra.CheckErr(err)
				_, err = output.Write(block)
				cobra.CheckErr(err)
				blocksWritten.Add(blocksWritten, tntengine.BigOne)
			}
		}

		last = r.Int64()
		for i := first; i < last; i++ {
			_, err := random.Read(block)
			cobra.CheckErr(err)
			_, err = output.Write(block)
			cobra.CheckErr(err)
			blocksWritten.Add(blocksWritten, tntengine.BigOne)
		}

		fmt.Fprintf(os.Stderr, "Blocks Written: %d\n", blocksWritten)
	}()

	return input
}

func aliasNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "bs":
		name = "bs=BYTES"
	case "count":
		name = "count=N"
	}
	return pflag.NormalizedName(name)
}

// checkFatal checks for error that are not io.EOF and io.ErrUnexpectedEOF and logs them.
func checkError(e error) {
	if e != io.EOF && e != io.ErrUnexpectedEOF {
		cobra.CheckErr(e)
	}
}

func readCounterFile(defaultMap map[string]*big.Int) map[string]*big.Int {
	f, err := os.OpenFile(cntrFileName, os.O_RDONLY, 0600)
	if err != nil {
		return defaultMap
	}

	defer f.Close()
	cmap := make(map[string]*big.Int)
	dec := gob.NewDecoder(f)
	checkError(dec.Decode(&cmap))
	return cmap
}

func writeCounterFile(wMap map[string]*big.Int) error {
	f, err := os.OpenFile(cntrFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer f.Close()
	enc := gob.NewEncoder(f)
	return enc.Encode(wMap)
}

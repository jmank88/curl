package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/exec"

	"github.com/posener/complete"

	"github.com/jmank88/curl/jsonrpc"
)

const env = "CURLETH"

var cfg jsonrpc.Config

func init() {
	flag.BoolVar(&cfg.Cmds, "x", false, "Print command")
	flag.BoolVar(&cfg.Nop, "n", false, "Print command without running")
	flag.BoolVar(&cfg.Pretty, "p", false, "Pretty JSON formatting")
	flag.BoolVar(&cfg.Verbose, "v", false, "Verbose logs")
	flag.IntVar(&cfg.ID, "id", -1, "JSONRPC ID - random if < 0")
}

var cmd = complete.Command{
	GlobalFlags: map[string]complete.Predictor{
		"-x":  complete.PredictNothing,
		"-n":  complete.PredictNothing,
		"-p":  complete.PredictNothing,
		"-v":  complete.PredictNothing,
		"-id": complete.PredictNothing,
	},
	Sub: map[string]complete.Command{
		"eth": {Sub: map[string]complete.Command{
			"accounts":                         {},
			"blockNumber":                      {},
			"coinbase":                         {},
			"compileLLL":                       {},
			"compileSerpent":                   {},
			"compileSolidity":                  {},
			"call":                             {},
			"estimateGas":                      {},
			"getBalance":                       {Args: predictor{predictHex, predictBool}},
			"getBlockTransactionCountByHash":   {Args: predictor{predictHex}},
			"getBlockTransactionCountByNumber": {Args: predictor{predictBlockNum}},
			"getBlockByHash":                   {Args: predictor{predictHex, predictBool}},
			"getBlockByNumber":                 {Args: predictor{predictBlockNum, predictBool}},
			"getCode":                          {Args: predictor{predictHex, predictBlockNum}},
			"getCompilers":                     {},
			"getFilterChanges":                 {Args: predictor{predictHex}},
			"getFilterLogs":                    {},
			"getLogs":                          {
				//TODO flags...from/to/address/topics/blockhash
			},
			"gasPrice":                             {},
			"getStorageAt":                         {Args: predictor{predictHex, predictHex, predictBlockNum}},
			"getTransactionCount":                  {Args: predictor{predictHex, predictBlockNum}},
			"getTransactionReceipt":                {Args: predictor{predictHex}},
			"getTransactionsByHash":                {Args: predictor{predictHex}},
			"getTransactionsByBlockHashAndIndex":   {Args: predictor{predictHex, predictHex}},
			"getTransactionsByBlockNumberAndIndex": {Args: predictor{predictBlockNum, predictHex}},
			"getUncleCountByBlockHash":             {Args: predictor{predictHex}},
			"getUncleCountByBlockNumber":           {Args: predictor{predictBlockNum}},
			"getUncleByBlockHashAndIndex":          {Args: predictor{predictHex, predictHex}},
			"getUncleByBlockNumberAndIndex":        {Args: predictor{predictBlockNum, predictHex}},
			"getWork":                              {Args: predictor{predictHex, predictHex, predictHex}},
			"hashrate":                             {},
			"mining":                               {},
			"newBlockFilter":                       {},
			"newFilter":                            {},
			"newPendingTransactionFilter":          {},
			"protocolVersion":                      {},
			"sign":                                 {Args: predictor{predictHex, predictHex}},
			"syncing":                              {},
			"signTransaction":                      {
				//TODO flags
			},
			"sendTransaction": {
				//TODO flags
			},
			"sendRawTransaction": {Args: predictor{predictHex}},
			"submitWork":         {Args: predictor{predictHex, predictHex, predictHex}},
			"submitHashrate":     {Args: predictor{predictHex, predictHex}},
			"uninstallFilter":    {},
		}},
		"net": {Sub: map[string]complete.Command{
			"listening": {},
			"peerCount": {},
			"version":   {},
		}},
		"web3": {Sub: map[string]complete.Command{
			"clientVersion": {},
			"sha3":          {},
		}},
	},
}

func main() {
	os.Exit(Main())
}

func Main() int {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintln(out, "Usage:")
		fmt.Fprintln(out, "  curleth (eth|net|web3) (method) [args...] [url]")
		fmt.Fprintln(out, "Environment:")
		fmt.Fprintln(out, "  CURLETH - set to omit url arg")
		fmt.Fprintln(out, "Flags:")
		flag.CommandLine.PrintDefaults()
	}
	if complete.New(os.Args[0], cmd).Run() {
		return 0 // autocompleted
	}

	need := 2
	url, envSet := os.LookupEnv(env)
	if !envSet {
		need++ // expect url as last arg instead
	}

	args := flag.CommandLine.Args()
	if l := len(args); l < need {
		s := "curleth (eth|net|web3) (method) [args...]"
		if !envSet {
			s += " [url]"
		}
		log.Printf("Too few arguments %d, need at least %d: %s\n", l, need, s)
		return 1
	}
	method := args[0] + "_" + args[1]
	args = args[2:]
	if !envSet {
		url = args[len(args)-1]
		args = args[:len(args)-1]
	}

	resp, err := cfg.Do(context.Background(), url, method, typedParams(args)...)
	if err != nil {
		if exit := new(exec.ExitError); errors.As(err, &exit) {
			log.Println(exit.String())
			return exit.ExitCode()
		}
		log.Println(err)
		return 1
	}
	fmt.Println(string(resp))
	return 0
}

func typedParams(params []string) []any {
	ts := make([]any, len(params))
	for i, p := range params {
		if p == "true" {
			ts[i] = true
		} else if p == "false" {
			ts[i] = false
		} else if v, ok := new(big.Int).SetString(p, 10); ok {
			ts[i] = "0x" + v.Text(16)
		} else {
			ts[i] = p
		}
	}
	return ts
}

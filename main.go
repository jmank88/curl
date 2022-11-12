package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/posener/complete"
)

const env = "ETHCURL"

var (
	verbose = flag.Bool("v", false, "verbose")
	nop     = flag.Bool("n", false, "print command without running")
	cmds    = flag.Bool("x", false, "print command that is run")

	cmd = complete.Command{
		GlobalFlags: map[string]complete.Predictor{
			"-v": complete.PredictNothing,
			"-n": complete.PredictNothing,
			"-x": complete.PredictNothing,
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
)

func main() {
	if complete.New(os.Args[0], cmd).Run() {
		return
	}

	need := 2
	url, envSet := os.LookupEnv(env)
	if !envSet {
		need++
	}

	args := flag.CommandLine.Args()
	if len(args) < need {
		log.Fatalln("Too few arguments: ethcurl (eth|net|web3) (method) [args...] [url]")
	}
	method := args[0] + "_" + args[1]
	args = args[2:]
	if !envSet {
		url = args[len(args)-1]
		args = args[:len(args)-1]
	}

	if err := execCurl(context.Background(), url, method, args...); err != nil {
		log.Fatalln(err)
	}
}

func execCurl(ctx context.Context, url, method string, params ...string) error {
	b, err := json.Marshal(struct {
		ID      string `json:"id"`
		Version string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  []any  `json:"params"`
	}{
		ID:      strconv.Itoa(rand.Int()),
		Version: "2.0",
		Method:  method,
		Params:  typedParams(params),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	var args []string
	if !*verbose {
		args = append(args, "-s") // --silent
	}
	args = append(args, "-H", "Content-Type: application/json", "-X", "POST", "-d", string(b), url)
	if *nop || *cmds || *verbose {
		fmt.Fprintln(os.Stderr, "curl", strings.Join(args, " "))
	}
	if *nop {
		return nil
	}

	cmd := exec.CommandContext(ctx, "curl", args...)
	cmd.Stderr = os.Stderr
	b, err = cmd.Output()
	if err != nil {
		if exit := new(exec.ExitError); errors.As(err, &exit) {
			log.Println(exit.String())
			os.Exit(exit.ExitCode())
		}
		return err
	}
	if *verbose {
		fmt.Fprint(os.Stderr, string(b))
	}
	resp := struct {
		ID      string          `json:"id"`
		Version string          `json:"jsonrpc"`
		Result  json.RawMessage `json:"result"`
		Error   *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}{}
	if err := json.Unmarshal(b, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("jsonrpc error: %w: %s", jsonrpcError(resp.Error.Code), resp.Error.Message)
	}

	fmt.Println(string(resp.Result))

	return nil
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

type jsonrpcError int

func (e jsonrpcError) Error() string {
	m := "Unrecognized error"
	switch e {
	case -32700:
		m = "parse error"
	case -32600:
		m = "invalid Request"
	case -32601:
		m = "method not found"
	case -32602:
		m = "invalid params"
	case -32603:
		m = "internal error"
	default:
		if -32000 > e && e > -32099 {
			m = "server error"
		}
	}

	return fmt.Sprintf("%d (%s)", e, m)
}

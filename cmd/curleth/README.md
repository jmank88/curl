# curleth

CLI to simplify Ethereum JSONRPC over `curl`.

## Features

- Autocomplete
- Hidden JSONRPC

## Install

`go install github.com/jmank88/curl/cmd/curleth`

### Bash Autocomplete (optional)

`curleth -install`

## How to use

`curleth (eth|net|web3) (method) [args...] [url]`

### Environment

- `CURLETH`: can be used instead of passing `url` as the last argument.

### Flags
- `-install`: Install completion for curleth command
- `-uninstall`: Uninstall completion for curleth command
- `-y`: Don't prompt user for typing 'yes' when installing completion
- `-n`: Print command without running
- `-v`: Verbose logs
- `-x`: Print command that is run

## Example

```bash
$ export CURLETH=https://rpc.gochain.io 

$ curleth eth getBlockByNumber latest false | jq
{
  "difficulty": "0x12",
  "extraData": "0x46696e7320417474616368656400000000000000000000000000000000000000",
  "gasLimit": "0x822d320",
  "gasUsed": "0x0",
  "hash": "0x71dcde4b54dd0965b277b415a1960a5074fe886dc5339c9cc2a6f6d6dd3b0c88",
  "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "miner": "0x8846359af59b723ab4c540456251b9c3fe2f269d",
  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "nonce": "0x0000000000000000",
  "number": "0x1b0947d",
  "parentHash": "0xecc7cc0ceb19608788a86b19dae827e3de22eea7b5bd5d1bff73245e067a772e",
  "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
  "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
  "signer": "0x3fc1029c1e90854fa357ffa0d19cf9184e754a12f5c7d8a627286d6b889dd8172de65cd2784ecd214d89303489f5b730805e7aebdd38fbff61c177d298d16c1600",
  "signers": [],
  "size": "0x266",
  "stateRoot": "0xfc728435863ac64cb3f353c6e46b1f30c0597605bb1e15c33d0203fb32d42c9c",
  "timestamp": "0x63701110",
  "totalDifficulty": "0x1cc955bd",
  "transactions": [],
  "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
  "uncles": [],
  "voters": []
}

$ curleth -x eth getBalance 0xEd7F2e81B0264177e0df8f275f97Fd74Fa51A896 latest
> curl -s -d {"id":"5577006791947779410","jsonrpc":"2.0","method":"eth_getBalance","params":["0xEd7F2e81B0264177e0df8f275f97Fd74Fa51A896","latest"]} -H Content-Type: application/json -X POST https://rpc.gochain.io
0x3da15d8c524996eb8eca

$ curleth eth getTransactionCount 0xEd7F2e81B0264177e0df8f275f97Fd74Fa51A896 latest
0xf
```

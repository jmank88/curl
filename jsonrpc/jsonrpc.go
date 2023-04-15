package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/jmank88/curl"
)

// Do POSTs a request and returns the raw result bytes, or an error.
// Responses which contain errors will be of the type *Error.
func Do(ctx context.Context, c curl.Config, url string, method string, params ...any) ([]byte, error) {
	var resp Response
	err := c.PostJSON(ctx, url, Request{Method: method, Params: params}, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

type Request struct {
	Version Version2 `json:"jsonrpc"`
	ID      RandID   `json:"id"`
	Method  string   `json:"method"`
	Params  []any    `json:"params"`
}

type Version2 string

func (Version2) MarshalJSON() ([]byte, error) { return []byte("2.0"), nil }

type RandID int

func (RandID) MarshalJSON() ([]byte, error) { return []byte(strconv.Itoa(rand.Int())), nil }

type Response struct {
	Version string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *Error          `json:"error"`
}

// Error represents a jsonrpc error, and formats with full details and a code description.
type Error struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// Error returns an error string of the form:
//
//	jsonrpc error: <code> (<description>): <message>[: data]
func (e *Error) Error() string {
	desc := "Unrecognized error"
	switch e.Code {
	case -32700:
		desc = "Parse error"
	case -32600:
		desc = "Invalid Request"
	case -32601:
		desc = "Method not found"
	case -32602:
		desc = "Invalid params"
	case -32603:
		desc = "Internal error"
	default:
		if -32000 > e.Code && e.Code > -32099 {
			desc = "Server error"
		}
	}

	s := fmt.Sprintf("jsonrpc error: %d (%s): %s", e.Code, desc, e.Message)
	if len(e.Data) == 0 {
		return s
	}
	return fmt.Sprintf("%s: %s", s, e.Data)
}

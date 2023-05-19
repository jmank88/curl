package curl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	Cmds    bool // Print command. Redundant with Nop or Verbose.
	Nop     bool // Print command without running.
	Pretty  bool // Pretty JSON formatting.
	Verbose bool // Verbose logs.

	Stderr io.Writer // Where to write error logs.
}

// Post calls Config.Post with the zero value defaults.
func Post(ctx context.Context, url string, data []byte, args ...string) ([]byte, error) {
	return Config{}.Post(ctx, url, data, args...)
}

// Post constructs a curl command to POST data to url using args, and executes it, unless Config.Nop is true.
func (c Config) Post(ctx context.Context, url string, data []byte, args ...string) ([]byte, error) {
	if c.Stderr == nil {
		c.Stderr = os.Stderr
	}

	if !c.Verbose {
		args = append(args, "-s") // --silent
	}
	args = append(args, "-X", "POST", "-d", string(data), url)
	if c.Nop || c.Cmds || c.Verbose {
		fmt.Fprintln(c.Stderr, "curl", strings.Join(args, " "))
	}
	if c.Nop {
		return nil, nil
	}

	cmd := exec.CommandContext(ctx, "curl", args...)
	cmd.Stderr = c.Stderr
	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if c.Verbose {
		fmt.Fprintln(c.Stderr, string(b))
	}
	return b, nil
}

// PostJSON calls Config.PostJSON with the zero value defaults.
func PostJSON(ctx context.Context, url string, req, resp any) error {
	return Config{}.PostJSON(ctx, url, req, resp)
}

// PostJSON calls Post with req marshalled as JSON, and unmarshalls in to resp.
func (c Config) PostJSON(ctx context.Context, url string, req, resp any) error {
	reqB, err := c.marshalJSON(req)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	respB, err := c.Post(ctx, url, reqB, "-H", "Content-Type: application/json")
	if err != nil {
		return fmt.Errorf("failed to post: %w", err)
	}
	if !c.Nop {
		if err := json.Unmarshal(respB, &resp); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w: %s", err, string(respB))
		}
	}
	return nil
}

func (c Config) marshalJSON(req any) ([]byte, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	if c.Pretty {
		e.SetIndent("", "  ")
	}
	err := e.Encode(req)
	if err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(b.Bytes(), []byte("\n")), nil
}

package curl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	Cmds    bool // Print command that is executed.
	Nop     bool // Print command without running.
	Verbose bool // Verbose logs.

	Stderr io.Writer
}

func Post(ctx context.Context, url string, data []byte, args ...string) ([]byte, error) {
	return Config{}.Post(ctx, url, data, args...)
}

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

// PostJSON calls Post with req marshalled as JSON, and unmarshalls in to resp.
func PostJSON(ctx context.Context, url string, req, resp any) error {
	return Config{}.PostJSON(ctx, url, req, resp)
}

func (c Config) PostJSON(ctx context.Context, url string, req, resp any) error {
	reqB, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}
	respB, err := c.Post(ctx, url, reqB, "-H", "Content-Type: application/json")
	if err := json.Unmarshal(respB, &resp); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}
	return nil
}
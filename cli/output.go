package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"wolfapi/api"
)

func printResponse(resp api.GenericResponse, raw []byte) error {
	cfg := ResolvedConfig()

	if cfg.Output == "json" {
		var buf bytes.Buffer
		if err := json.Indent(&buf, raw, "", "  "); err != nil {
			_, writeErr := os.Stdout.Write(raw)
			if writeErr != nil {
				return writeErr
			}
			_, writeErr = fmt.Fprintln(os.Stdout)
			return writeErr
		}
		_, err := fmt.Fprintln(os.Stdout, buf.String())
		return err
	}

	if resp.Result != "" {
		_, err := fmt.Fprintln(os.Stdout, resp.Result)
		return err
	}

	_, err := os.Stdout.Write(raw)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(os.Stdout)
	return err
}

func printBatchHeader(label string) error {
	_, err := fmt.Fprintf(os.Stdout, "== %s ==\n", label)
	return err
}

func printBatchError(err error) error {
	_, writeErr := fmt.Fprintf(os.Stderr, "error: %v\n", err)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

package cmd

import (
	"bytes"
	"fmt"
	"os"
)

type promptConfig struct {
	raw  string
	file string
}

func (p promptConfig) resolve() (string, error) {
	if p.raw != "" {
		return p.raw, nil
	}

	if p.file == "" {
		return "", fmt.Errorf("prompt required (use --prompt or --prompt-file)")
	}

	data, err := os.ReadFile(p.file)
	if err != nil {
		return "", fmt.Errorf("read prompt file %q: %w", p.file, err)
	}

	return string(bytes.TrimSpace(data)), nil
}

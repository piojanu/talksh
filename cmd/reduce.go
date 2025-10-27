package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/piojanu/talksh/pkg/llm"
	"github.com/spf13/cobra"
)

var reduceCmd = &cobra.Command{
	Use:   "reduce",
	Short: "Run a prompt once on the entire stdin payload",
	Long: `Run a prompt once on the entire stdin payload

Whole-input design: you one the upstream aggregation.
hatever you feed us (via cat, jq, etc.) becomes on LLM call.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt, err := reducePrompt.resolve()
		if err != nil {
			return err
		}

		payload, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}

		resp, err := llm.CallLLM(prompt, string(bytes.TrimSpace(payload)))
		if err != nil {
			return fmt.Errorf("call LLM: %w", err)
		}

		fmt.Fprintln(os.Stdout, resp)
		return nil
	},
}

var (
	reducePrompt promptConfig
)

func init() {
	reduceCmd.Flags().StringVar(&reducePrompt.raw, "prompt", "", "raw prompt string")
	reduceCmd.Flags().StringVar(&reducePrompt.file, "prompt-file", "", "prompt file path")
	rootCmd.AddCommand(reduceCmd)
}

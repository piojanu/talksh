package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/piojanu/talksh/pkg/llm"
	"github.com/spf13/cobra"
)

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Run a prompt on each input line",
	Long: `Run a prompt on each input line

Intentional line-by-line design: one stdin line == one LLM call.
You must make upstream tools (cat, jq, etc.) emit exactly what you want per line.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt, err := mapPrompt.resolve()
		if err != nil {
			return err
		}

		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			line := scanner.Text()

			resp, err := llm.CallLLM(prompt, line)
			if err != nil {
				return fmt.Errorf("call LLM: %w", err)
			}

			fmt.Fprintln(os.Stdout, resp)
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("read stdin: %w", err)
		}

		return nil
	},
}

var (
	mapPrompt promptConfig
)

func init() {
	mapCmd.Flags().StringVar(&mapPrompt.raw, "prompt", "", "raw prompt string")
	mapCmd.Flags().StringVar(&mapPrompt.file, "prompt-file", "", "prompt file path")
	rootCmd.AddCommand(mapCmd)
}

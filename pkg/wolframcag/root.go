package wolframcag

import (
	"fmt"

	"github.com/spf13/cobra"
)

var flags Config

var resolvedConfig Config
var resolvedClient *Client

var rootCmd = &cobra.Command{
	Use:   "wolfram-cag",
	Short: "CLI for Wolfram CAG APIs",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := Resolve(flags)
		if err != nil {
			return err
		}

		if err := validateOutputFormat(cfg.Output); err != nil {
			return err
		}

		resolvedConfig = cfg
		resolvedClient = New(cfg)

		if flags.Verbose {
			fmt.Printf("resolved config: base-url=%s output=%s timeout-secs=%d workers=%d\n", cfg.BaseURL, cfg.Output, cfg.TimeoutSecs, cfg.Workers)
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flags.APIKey, "api-key", "", "Wolfram API key (overrides WOLFRAM_APP_ID)")
	rootCmd.PersistentFlags().StringVar(&flags.BaseURL, "base-url", DefaultBaseURL, "Wolfram CAG base URL")
	rootCmd.PersistentFlags().StringVar(&flags.Output, "output", DefaultOutput, "Output format: text|json")
	rootCmd.PersistentFlags().IntVar(&flags.TimeoutSecs, "timeout-secs", DefaultTimeout, "HTTP timeout in seconds")
	rootCmd.PersistentFlags().BoolVar(&flags.Verbose, "verbose", false, "Enable verbose logging")
	rootCmd.PersistentFlags().IntVar(&flags.Workers, "workers", 4, "Number of worker goroutines for batch operations")

	rootCmd.AddCommand(ctxAPICmd)
	rootCmd.AddCommand(resultCmd)
	rootCmd.AddCommand(computeCmd)
	rootCmd.AddCommand(hintsCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func ResolvedConfig() Config {
	return resolvedConfig
}

func ResolvedClient() *Client {
	return resolvedClient
}

func validateOutputFormat(output string) error {
	if output == "text" || output == "json" {
		return nil
	}
	return fmt.Errorf("invalid --output %q: must be one of text,json", output)
}

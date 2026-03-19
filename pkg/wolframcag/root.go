package wolframcag

import (
	"fmt"

	"github.com/spf13/cobra"
)

type globalFlags struct {
	apiKey      string
	baseURL     string
	output      string
	timeoutSecs int
	verbose     bool
	workers     int
}

var flags globalFlags

var resolvedConfig Config
var resolvedClient *Client

var rootCmd = &cobra.Command{
	Use:   "wolfram-cag",
	Short: "CLI for Wolfram CAG APIs",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := Resolve(Inputs{
			APIKey:      flags.apiKey,
			BaseURL:     flags.baseURL,
			Output:      flags.output,
			TimeoutSecs: flags.timeoutSecs,
			Verbose:     flags.verbose,
			Workers:     flags.workers,
		})
		if err != nil {
			return err
		}

		if err := validateOutputFormat(cfg.Output); err != nil {
			return err
		}

		resolvedConfig = cfg
		resolvedClient = New(cfg)

		if flags.verbose {
			fmt.Printf("resolved config: base-url=%s output=%s timeout-secs=%d workers=%d\n", cfg.BaseURL, cfg.Output, cfg.TimeoutSecs, cfg.Workers)
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flags.apiKey, "api-key", "", "Wolfram API key (overrides WOLFRAM_APP_ID)")
	rootCmd.PersistentFlags().StringVar(&flags.baseURL, "base-url", DefaultBaseURL, "Wolfram CAG base URL")
	rootCmd.PersistentFlags().StringVar(&flags.output, "output", DefaultOutput, "Output format: text|json")
	rootCmd.PersistentFlags().IntVar(&flags.timeoutSecs, "timeout-secs", DefaultTimeout, "HTTP timeout in seconds")
	rootCmd.PersistentFlags().BoolVar(&flags.verbose, "verbose", false, "Enable verbose logging")
	rootCmd.PersistentFlags().IntVar(&flags.workers, "workers", 4, "Number of worker goroutines for batch operations")

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

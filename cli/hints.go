package cli

import "github.com/spf13/cobra"

var hintsContext string

var hintsCmd = &cobra.Command{
	Use:   "hints",
	Short: "Call WolframLanguageHints API",
	RunE: func(cmd *cobra.Command, args []string) error {
		if hintsContext == "" {
			return cmd.Usage()
		}
		return nil
	},
}

func init() {
	hintsCmd.Flags().StringVar(&hintsContext, "context", "", "Prompt or context text")
	_ = hintsCmd.MarkFlagRequired("context")
}

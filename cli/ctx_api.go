package cli

import "github.com/spf13/cobra"

var ctxAPIText string

var ctxAPICmd = &cobra.Command{
	Use:   "context",
	Short: "Call WolframAlphaContext API",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ctxAPIText == "" {
			return cmd.Usage()
		}
		return nil
	},
}

func init() {
	ctxAPICmd.Flags().StringVar(&ctxAPIText, "context", "", "Conversation context text")
	_ = ctxAPICmd.MarkFlagRequired("context")
}

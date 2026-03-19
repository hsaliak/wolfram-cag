package cli

import (
	"github.com/spf13/cobra"

	"wolfapi/api"
)

var ctxAPIText string

var ctxAPICmd = &cobra.Command{
	Use:   "context",
	Short: "Call WolframAlphaContext API",
	RunE: func(cmd *cobra.Command, args []string) error {
		if ctxAPIText == "" {
			return cmd.Usage()
		}

		svc := api.New(ResolvedClient())
		resp, raw, err := svc.Context(cmd.Context(), api.ContextRequest{Context: ctxAPIText})
		if err != nil {
			return err
		}

		return printResponse(resp, raw)
	},
}

func init() {
	ctxAPICmd.Flags().StringVar(&ctxAPIText, "context", "", "Conversation context text")
	_ = ctxAPICmd.MarkFlagRequired("context")
}

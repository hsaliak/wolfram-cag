package cli

import (
	"github.com/spf13/cobra"

	"wolfapi/api"
)

var hintsContext string

var hintsCmd = &cobra.Command{
	Use:   "hints",
	Short: "Call WolframLanguageHints API",
	RunE: func(cmd *cobra.Command, args []string) error {
		if hintsContext == "" {
			return cmd.Usage()
		}

		svc := api.New(ResolvedClient())
		resp, raw, err := svc.Hints(cmd.Context(), api.HintsRequest{Context: hintsContext})
		if err != nil {
			return err
		}

		return printResponse(resp, raw)
	},
}

func init() {
	hintsCmd.Flags().StringVar(&hintsContext, "context", "", "Prompt or context text")
	_ = hintsCmd.MarkFlagRequired("context")
}

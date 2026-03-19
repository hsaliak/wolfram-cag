package cli

import (
	"github.com/spf13/cobra"

	"wolfapi/api"
)

var hintsCmd = &cobra.Command{
	Use:   "hints <context-text>",
	Short: "Call WolframLanguageHints API",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hintsText := args[0]

		svc := api.New(ResolvedClient())
		resp, raw, err := svc.Hints(cmd.Context(), api.HintsRequest{Context: hintsText})
		if err != nil {
			return err
		}

		return printResponse(resp, raw)
	},
}

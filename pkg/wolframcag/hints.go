package wolframcag

import "github.com/spf13/cobra"

var hintsCmd = &cobra.Command{
	Use:   "hints <context-text>",
	Short: "Call WolframLanguageHints API",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hintsText := args[0]

		svc := NewService(ResolvedClient())
		resp, raw, err := svc.Hints(cmd.Context(), HintsRequest{Context: hintsText})
		if err != nil {
			return err
		}

		return printResponse(resp, raw)
	},
}

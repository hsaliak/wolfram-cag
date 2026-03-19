package wolframcag

import "github.com/spf13/cobra"

var ctxAPICmd = &cobra.Command{
	Use:   "context <context-text>",
	Short: "Call WolframAlphaContext API",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		contextText := args[0]

		svc := NewService(ResolvedClient())
		resp, raw, err := svc.Context(cmd.Context(), ContextRequest{Context: contextText})
		if err != nil {
			return err
		}

		return printResponse(resp, raw)
	},
}

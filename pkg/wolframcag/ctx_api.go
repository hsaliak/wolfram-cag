package wolframcag

import (
	"context"

	"github.com/spf13/cobra"
)

var ctxAPICmd = &cobra.Command{
	Use:   "context <context-text>",
	Short: "Call WolframAlphaContext API",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := NewService(ResolvedClient())
		return runSingleTextCommand(cmd, args, func(ctx context.Context, text string) (GenericResponse, []byte, error) {
			return svc.Context(ctx, ContextRequest{Context: text})
		})
	},
}

func runSingleTextCommand(cmd *cobra.Command, args []string, call func(context.Context, string) (GenericResponse, []byte, error)) error {
	resp, raw, err := call(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	return printResponse(resp, raw)
}

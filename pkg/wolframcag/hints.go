package wolframcag

import (
	"context"

	"github.com/spf13/cobra"
)

var hintsCmd = &cobra.Command{
	Use:   "hints <context-text>",
	Short: "Call WolframLanguageHints API",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := NewService(ResolvedClient())
		return runSingleTextCommand(cmd, args, func(ctx context.Context, text string) (GenericResponse, []byte, error) {
			return svc.Hints(ctx, HintsRequest{Context: text})
		})
	},
}

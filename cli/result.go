package cli

import "github.com/spf13/cobra"

var (
	resultInput     string
	resultInputFile string
)

var resultCmd = &cobra.Command{
	Use:   "result",
	Short: "Call WolframAlphaResult API",
	RunE: func(cmd *cobra.Command, args []string) error {
		if resultInput == "" && resultInputFile == "" {
			return cmd.Usage()
		}
		return nil
	},
}

func init() {
	resultCmd.Flags().StringVar(&resultInput, "input", "", "Input query string")
	resultCmd.Flags().StringVar(&resultInputFile, "input-file", "", "Path to file with newline-delimited input queries")

	resultCmd.MarkFlagsMutuallyExclusive("input", "input-file")
	resultCmd.MarkFlagsOneRequired("input", "input-file")
}

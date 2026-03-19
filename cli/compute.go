package cli

import "github.com/spf13/cobra"

var (
	computeCode     string
	computeCodeFile string
)

var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Call WolframLanguageCompute API",
	RunE: func(cmd *cobra.Command, args []string) error {
		if computeCode == "" && computeCodeFile == "" {
			return cmd.Usage()
		}
		return nil
	},
}

func init() {
	computeCmd.Flags().StringVar(&computeCode, "code", "", "Wolfram Language code")
	computeCmd.Flags().StringVar(&computeCodeFile, "code-file", "", "Path to file with newline-delimited Wolfram Language code")

	computeCmd.MarkFlagsMutuallyExclusive("code", "code-file")
	computeCmd.MarkFlagsOneRequired("code", "code-file")
}

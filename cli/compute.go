package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"wolfapi/api"
)

var (
	computeCode           string
	computeCodeFile       string
	computeTimeConstraint int
	computeLine           int
	computeMaxChars       int
)

var computeCmd = &cobra.Command{
	Use:   "compute",
	Short: "Call WolframLanguageCompute API",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := api.New(ResolvedClient())
		opts := api.ComputeOptions{TimeConstraint: computeTimeConstraint, Line: computeLine, MaxChars: computeMaxChars}

		if computeCodeFile == "" {
			resp, raw, err := svc.Compute(cmd.Context(), api.ComputeRequest{Code: computeCode}, opts)
			if err != nil {
				return err
			}
			return printResponse(resp, raw)
		}

		inputs, err := readNonEmptyLines(computeCodeFile)
		if err != nil {
			return err
		}

		results := runStringBatch(inputs, ResolvedConfig().Workers, func(in string) batchResult {
			resp, raw, callErr := svc.Compute(cmd.Context(), api.ComputeRequest{Code: in}, opts)
			return batchResult{label: in, resp: resp, raw: raw, err: callErr}
		})

		hadErr := false
		for _, item := range results {
			_ = printBatchHeader(item.label)
			if item.err != nil {
				hadErr = true
				_ = printBatchError(item.label, item.err)
				continue
			}
			if err := printResponse(item.resp, item.raw); err != nil {
				return err
			}
		}

		if hadErr {
			return fmt.Errorf("compute batch completed with errors")
		}
		return nil
	},
}

func init() {
	computeCmd.Flags().StringVar(&computeCode, "code", "", "Wolfram Language code")
	computeCmd.Flags().StringVar(&computeCodeFile, "code-file", "", "Path to file with newline-delimited Wolfram Language code")

	computeCmd.MarkFlagsMutuallyExclusive("code", "code-file")
	computeCmd.MarkFlagsOneRequired("code", "code-file")
	computeCmd.Flags().IntVar(&computeTimeConstraint, "time-constraint", 0, "Time constraint for compute request")
	computeCmd.Flags().IntVar(&computeLine, "line", 0, "Line selection for compute request")
	computeCmd.Flags().IntVar(&computeMaxChars, "max-chars", 0, "Max chars for compute request")
}

package wolframcag

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	computeCodeFile       string
	computeTimeConstraint int
	computeLine           int
	computeMaxChars       int
)

func resolveComputeInput(args []string, codeFile string) (singleCode string, useBatch bool, err error) {
	if codeFile != "" {
		if len(args) != 0 {
			return "", false, fmt.Errorf("compute accepts either positional code or --code-file, not both")
		}
		return "", true, nil
	}

	if len(args) != 1 {
		return "", false, fmt.Errorf("compute requires exactly one positional code or --code-file")
	}

	return args[0], false, nil
}

var computeCmd = &cobra.Command{
	Use:   "compute <code>",
	Short: "Call WolframLanguageCompute API",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := NewService(ResolvedClient())
		opts := ComputeOptions{TimeConstraint: computeTimeConstraint, Line: computeLine, MaxChars: computeMaxChars}

		singleCode, useBatch, err := resolveComputeInput(args, computeCodeFile)
		if err != nil {
			return err
		}

		if !useBatch {
			resp, raw, err := svc.Compute(cmd.Context(), ComputeRequest{Code: singleCode}, opts)
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
			resp, raw, callErr := svc.Compute(cmd.Context(), ComputeRequest{Code: in}, opts)
			return batchResult{label: in, resp: resp, raw: raw, err: callErr}
		})

		hadErr := false
		for _, item := range results {
			_ = printBatchHeader(item.label)
			if item.err != nil {
				hadErr = true
				_ = printBatchError(item.err)
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
	computeCmd.Flags().StringVar(&computeCodeFile, "code-file", "", "Path to file with newline-delimited Wolfram Language code")
	computeCmd.Flags().IntVar(&computeTimeConstraint, "time-constraint", 0, "Time constraint for compute request")
	computeCmd.Flags().IntVar(&computeLine, "line", 0, "Line selection for compute request")
	computeCmd.Flags().IntVar(&computeMaxChars, "max-chars", 0, "Max chars for compute request")
}

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"wolfapi/api"
)

var (
	resultInputFile  string
	resultAssumption string
	resultFormat     string
	resultUnits      string
	resultLocation   string
	resultLatLong    string
	resultTimeout    string
	resultMaxWidth   string
)

func resolveResultInput(args []string, inputFile string) (singleInput string, useBatch bool, err error) {
	if inputFile != "" {
		if len(args) != 0 {
			return "", false, fmt.Errorf("result accepts either positional input or --input-file, not both")
		}
		return "", true, nil
	}

	if len(args) != 1 {
		return "", false, fmt.Errorf("result requires exactly one positional input or --input-file")
	}

	return args[0], false, nil
}

var resultCmd = &cobra.Command{
	Use:   "result <input>",
	Short: "Call WolframAlphaResult API",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := api.New(ResolvedClient())
		opts := api.ResultOptions{
			Assumption: resultAssumption,
			Format:     resultFormat,
			Units:      resultUnits,
			Location:   resultLocation,
			LatLong:    resultLatLong,
			Timeout:    resultTimeout,
			MaxWidth:   resultMaxWidth,
		}

		singleInput, useBatch, err := resolveResultInput(args, resultInputFile)
		if err != nil {
			return err
		}

		if !useBatch {
			resp, raw, err := svc.Result(cmd.Context(), singleInput, opts)
			if err != nil {
				return err
			}
			return printResponse(resp, raw)
		}

		inputs, err := readNonEmptyLines(resultInputFile)
		if err != nil {
			return err
		}

		results := runStringBatch(inputs, ResolvedConfig().Workers, func(in string) batchResult {
			resp, raw, callErr := svc.Result(cmd.Context(), in, opts)
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
			return fmt.Errorf("result batch completed with errors")
		}
		return nil
	},
}

func init() {
	resultCmd.Flags().StringVar(&resultInputFile, "input-file", "", "Path to file with newline-delimited input queries")
	resultCmd.Flags().StringVar(&resultAssumption, "assumption", "", "Wolfram assumption string")
	resultCmd.Flags().StringVar(&resultFormat, "format", "", "Result format hint")
	resultCmd.Flags().StringVar(&resultUnits, "units", "", "Units preference")
	resultCmd.Flags().StringVar(&resultLocation, "location", "", "Location hint")
	resultCmd.Flags().StringVar(&resultLatLong, "latlong", "", "Latitude/longitude hint")
	resultCmd.Flags().StringVar(&resultTimeout, "timeout", "", "API timeout parameter")
	resultCmd.Flags().StringVar(&resultMaxWidth, "maxwidth", "", "Max width parameter")
}

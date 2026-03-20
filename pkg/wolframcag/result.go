package wolframcag

import (
	"fmt"

	"github.com/spf13/cobra"
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

var resultCmd = &cobra.Command{
	Use:   "result <input>",
	Short: "Call WolframAlphaResult API",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := NewService(ResolvedClient())
		opts := ResultOptions{
			Assumption: resultAssumption,
			Format:     resultFormat,
			Units:      resultUnits,
			Location:   resultLocation,
			LatLong:    resultLatLong,
			Timeout:    resultTimeout,
			MaxWidth:   resultMaxWidth,
		}

		singleInput, useBatch, err := resolveSingleArgOrFile(args, resultInputFile, "result", "input", "input-file")
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

		hadErr, err := printBatchResults(results)
		if err != nil {
			return err
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

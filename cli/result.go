package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"wolfapi/api"
)

var (
	resultInput      string
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
	Use:   "result",
	Short: "Call WolframAlphaResult API",
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

		if resultInputFile == "" {
			resp, raw, err := svc.Result(cmd.Context(), resultInput, opts)
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
				_ = printBatchError(item.label, item.err)
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
	resultCmd.Flags().StringVar(&resultInput, "input", "", "Input query string")
	resultCmd.Flags().StringVar(&resultInputFile, "input-file", "", "Path to file with newline-delimited input queries")

	resultCmd.MarkFlagsMutuallyExclusive("input", "input-file")
	resultCmd.MarkFlagsOneRequired("input", "input-file")
	resultCmd.Flags().StringVar(&resultAssumption, "assumption", "", "Wolfram assumption string")
	resultCmd.Flags().StringVar(&resultFormat, "format", "", "Result format hint")
	resultCmd.Flags().StringVar(&resultUnits, "units", "", "Units preference")
	resultCmd.Flags().StringVar(&resultLocation, "location", "", "Location hint")
	resultCmd.Flags().StringVar(&resultLatLong, "latlong", "", "Latitude/longitude hint")
	resultCmd.Flags().StringVar(&resultTimeout, "timeout", "", "API timeout parameter")
	resultCmd.Flags().StringVar(&resultMaxWidth, "maxwidth", "", "Max width parameter")
}

package wolframcag

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type batchResult struct {
	label string
	resp  GenericResponse
	raw   []byte
	err   error
}

type indexedInput struct {
	idx int
	val string
}

type indexedResult struct {
	idx int
	res batchResult
}

func readNonEmptyLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no non-empty lines found in %s", path)
	}

	return out, nil
}

func runStringBatch(inputs []string, workers int, fn func(string) batchResult) []batchResult {
	if workers <= 0 {
		workers = 1
	}

	jobs := make(chan indexedInput)
	results := make(chan indexedResult)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				results <- indexedResult{idx: job.idx, res: fn(job.val)}
			}
		}()
	}

	go func() {
		for i, in := range inputs {
			jobs <- indexedInput{idx: i, val: in}
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	ordered := make([]batchResult, len(inputs))
	for item := range results {
		ordered[item.idx] = item.res
	}

	return ordered
}

func resolveSingleArgOrFile(args []string, filePath, commandName, argName, fileFlag string) (single string, useBatch bool, err error) {
	if filePath != "" {
		if len(args) != 0 {
			return "", false, fmt.Errorf("%s accepts either positional %s or --%s, not both", commandName, argName, fileFlag)
		}
		return "", true, nil
	}

	if len(args) != 1 {
		return "", false, fmt.Errorf("%s requires exactly one positional %s or --%s", commandName, argName, fileFlag)
	}

	return args[0], false, nil
}

func printBatchResults(results []batchResult) (hadErr bool, err error) {
	for _, item := range results {
		_ = printBatchHeader(item.label)
		if item.err != nil {
			hadErr = true
			_ = printBatchError(item.err)
			continue
		}
		if err := printResponse(item.resp, item.raw); err != nil {
			return hadErr, err
		}
	}

	return hadErr, nil
}

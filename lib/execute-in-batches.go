package lib

import (
    "sync"
    "fmt"
)

func ExecuteInBatches[TResult any, TInput any](inputs []TInput, batchSize int, fn func(TInput) TResult) []TResult {
    var wg sync.WaitGroup
    results := make([]TResult, 0, len(inputs))
    fmt.Printf("Executing %d in batches of %d\n", len(inputs), batchSize)

    for i := 0; i < len(inputs); i += batchSize {
        end := i + batchSize
        if end > len(inputs) {
            end = len(inputs)
        }

        wg.Add(end - i)
        fmt.Printf("Adding %d requests\n", end - i)

        for j := i; j < end; j++ {
            fmt.Printf("Firing request number %d\n", j)
            go func(input TInput) {
                defer wg.Done()
                res := fn(input)
                results = append(results, res)
            }(inputs[j])
        }
        wg.Wait()
    }

    return results
}
package awsservice

import (
	"fmt"

	types2 "github.com/cduggn/ccexplorer/internal/types"
)

// maxPages bounds pagination. Cost Explorer bills per request, so a runaway
// token loop is expensive as well as slow.
const maxPages = 100

// paginate drives a token-based Cost Explorer list operation. fetch requests
// one page for the given token (nil on the first call); nextToken reads the
// token to follow from a page; merge folds a later page into the first
// page's result (the accumulator), which is mutated and returned. Pagination
// stops once nextToken reports no more pages, or after maxPages as a
// safety bound.
func paginate[TOut any](
	fetch func(nextToken *string) (TOut, error),
	nextToken func(TOut) *string,
	merge func(acc, page TOut) TOut,
) (TOut, error) {

	var result TOut
	var token *string

	for page := 1; ; page++ {
		out, err := fetch(token)
		if err != nil {
			return result, types2.APIError{
				Msg: err.Error(),
			}
		}

		if page == 1 {
			result = out
		} else {
			result = merge(result, out)
		}

		next := nextToken(out)
		if next == nil || *next == "" {
			break
		}

		if page >= maxPages {
			return result, types2.APIError{
				Msg: fmt.Sprintf(
					"result set exceeded %d pages; narrow the date range, "+
						"granularity or grouping", maxPages),
			}
		}

		token = next
	}

	return result, nil
}

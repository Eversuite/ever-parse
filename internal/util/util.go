package util

import "fmt"

// Check if the error exists (err != nil). If so the error and relating data is printed and the program panics.
// Params:
// err: error - The error itself
// data: ...any - all data objects to print
func Check(err error, data ...any) {
	if err == nil {
		return
	}
	fmt.Printf("[Critical] Error{%s} Data [%+v\n]", err, data)
	panic(err)

}
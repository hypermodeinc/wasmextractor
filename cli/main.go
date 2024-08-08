/*
 * Copyright 2024 Hypermode, Inc.
 * MIT License - https://opensource.org/licenses/MIT
 */

package main

import (
	"fmt"
	"os"

	"github.com/hypermodeAI/wasmextractor"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <wasm file>\n", os.Args[0])
		return
	}

	wasmFilePath := os.Args[1]

	wasmBytes, err := wasmextractor.ReadWasmFile(wasmFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	info, err := wasmextractor.ExtractWasmInfo(wasmBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println()

	fmt.Println("IMPORTS")
	fmt.Println("-------")
	if len(info.Imports) == 0 {
		fmt.Println("(none)")
	} else {
		for _, item := range info.Imports {
			fmt.Println(&item)
		}
	}
	fmt.Println()

	fmt.Println("EXPORTS")
	fmt.Println("-------")
	if len(info.Exports) == 0 {
		fmt.Println("(none)")
	} else {
		for _, item := range info.Exports {
			fmt.Println(&item)
		}
	}
	fmt.Println()
}

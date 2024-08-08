/*
 * Copyright 2024 Hypermode, Inc.
 * MIT License - https://opensource.org/licenses/MIT
 */

package wasmextractor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type WasmInfo struct {
	Imports []WasmItem
	Exports []WasmItem
}

type WasmItem struct {
	Name  string
	Kind  WasmItemKind
	Index uint32
}

func (i *WasmItem) String() string {
	return fmt.Sprintf("%s: %s (index: %d)", i.Kind, i.Name, i.Index)
}

type WasmItemKind int32

const (
	WasmFunction WasmItemKind = 0
	WasmTable    WasmItemKind = 1
	WasmMemory   WasmItemKind = 2
	WasmGlobal   WasmItemKind = 3
)

func (k WasmItemKind) String() string {
	switch k {
	case WasmFunction:
		return "Function"
	case WasmTable:
		return "Table"
	case WasmMemory:
		return "Memory"
	case WasmGlobal:
		return "Global"
	default:
		return "Unknown"
	}
}

func ReadWasmFile(wasmFilePath string) ([]byte, error) {
	wasmBytes, err := os.ReadFile(wasmFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading wasm file: %v", err)
	}

	magic := []byte{0x00, 0x61, 0x73, 0x6D} // "\0asm"
	if len(wasmBytes) < 8 || !bytes.Equal(wasmBytes[:4], magic) {
		return nil, fmt.Errorf("invalid wasm file")
	}

	if binary.LittleEndian.Uint32(wasmBytes[4:8]) != 1 {
		return nil, fmt.Errorf("unsupported wasm version")
	}

	return wasmBytes, nil
}

func ExtractWasmInfo(wasmBytes []byte) (*WasmInfo, error) {
	info := &WasmInfo{}

	offset := 8
	for offset < len(wasmBytes) {
		sectionID := wasmBytes[offset]
		offset++

		size, n := binary.Uvarint(wasmBytes[offset:])
		offset += n

		switch sectionID {
		case 2: // Import section
			info.Imports = readImports(wasmBytes[offset : offset+int(size)])

		case 7: // Export section
			info.Exports = readExports(wasmBytes[offset : offset+int(size)])
		}

		offset += int(size)
	}

	return info, nil
}

func readImports(data []byte) []WasmItem {

	numItems, n := binary.Uvarint(data)
	offset := n

	imports := make([]WasmItem, numItems)

	for i := 0; i < int(numItems); i++ {
		moduleLen, n := binary.Uvarint(data[offset:])
		offset += n

		moduleName := string(data[offset : offset+int(moduleLen)])
		offset += int(moduleLen)

		fieldLen, n := binary.Uvarint(data[offset:])
		offset += n

		fieldName := string(data[offset : offset+int(fieldLen)])
		offset += int(fieldLen)

		kind := data[offset]
		offset++

		index, n := binary.Uvarint(data[offset:])
		offset += n

		imports[i] = WasmItem{
			Name:  fmt.Sprintf("%s.%s", moduleName, fieldName),
			Kind:  WasmItemKind(kind),
			Index: uint32(index),
		}
	}

	return imports
}

func readExports(data []byte) []WasmItem {

	numItems, n := binary.Uvarint(data)
	offset := n

	exports := make([]WasmItem, numItems)

	for i := 0; i < int(numItems); i++ {
		fieldLen, n := binary.Uvarint(data[offset:])
		offset += n

		fieldName := string(data[offset : offset+int(fieldLen)])
		offset += int(fieldLen)

		kind := data[offset]
		offset++

		index, n := binary.Uvarint(data[offset:])
		offset += n

		exports[i] = WasmItem{
			Name:  fieldName,
			Kind:  WasmItemKind(kind),
			Index: uint32(index),
		}
	}

	return exports
}
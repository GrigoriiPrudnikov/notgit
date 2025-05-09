package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrintStruct(v any) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Println("PrintStruct: failed to encode:", err)
		return
	}
	fmt.Print(buf.String())
}

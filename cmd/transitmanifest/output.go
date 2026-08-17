package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func writeJSON(writer io.Writer, value any) error {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(writer, string(encoded))
	return err
}

func usage(writer io.Writer) {
	fmt.Fprintln(writer, "transitmanifest commands:")
	fmt.Fprintln(writer, "  create --id ID --trip ID --plate PLATE --origin NODE --destination NODE")
	fmt.Fprintln(writer, "  scan --manifest ID --scan ID --package ID --station NODE --operation load|unload")
	fmt.Fprintln(writer, "  close --manifest ID [--reason TEXT]")
	fmt.Fprintln(writer, "  handoff --manifest ID --from NODE --to NODE --signer ID")
	fmt.Fprintln(writer, "  list | report --date YYYY-MM-DD | audit --manifest ID | verify")
}

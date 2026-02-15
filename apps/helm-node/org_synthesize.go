package main

import (
	"fmt"
	"io"
)

func runOrgSynthesize(args []string, w io.Writer) {
	_, _ = fmt.Fprintln(w, "Organization Synthesis requires enterprise components not present in this build.")
}

package main

import (
	"fmt"
	"strings"

	"github.com/kylelemons/godebug/pretty"
)

var (
	// pc is the pretty.Config used by the pp function.
	pc = &pretty.Config{
		IncludeUnexported: true,
	}

	// pp is an alias to pretty.Sprint with the pc config.
	pp = pc.Sprint
)

func errorContains(gotErr, wantErr error) bool {
	return strings.Contains(fmt.Sprint(gotErr), fmt.Sprint(wantErr))
}

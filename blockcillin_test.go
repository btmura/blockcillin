package main

import (
	"fmt"
	"strings"
)

func errorContains(gotErr, wantErr error) bool {
	return strings.Contains(fmt.Sprintf("%v", gotErr), fmt.Sprintf("%v", wantErr))
}

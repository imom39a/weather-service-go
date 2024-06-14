//go:build tools
// +build tools

package tools

/*

This package is used to import the tools as blank imports so that they
are available in the go.mod file but not in the codebase.

The tools are used for various purposes like code generation, linting, etc. The
tools are not used in the production code, but they are used in the development
process to automate tasks and improve the code quality.

*/

import (
	_ "github.com/deepmap/oapi-codegen/cmd/oapi-codegen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)

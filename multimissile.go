package msl

import (
	"fmt"
	"runtime"

	"github.com/istyle-inc/multimissile/config"
	"github.com/istyle-inc/multimissile/wlog"
)

const (
	Version = "0.3.1"
)

var (
	Config config.Config
	AL     wlog.Logger
	EL     wlog.Logger
)

func ServerHeader() string {
	return fmt.Sprintf("multimissile %s", Version)
}

func PrintVersion() {
	fmt.Printf(`msl %s
Compiler: %s %s
Copyright (C) 2016 Mercari, Inc.
`,
		Version,
		runtime.Compiler,
		runtime.Version())
}

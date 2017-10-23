package msl

import (
	"fmt"
	"runtime"

	"github.com/istyle-inc/multimissile/config"
	"github.com/istyle-inc/multimissile/wlog"
)

const (
	Version    = "0.0.1"
	WbtVersion = "0.3.1"
)

var (
	Config config.Config
	AL     wlog.Logger
	EL     wlog.Logger
)

func ServerHeader() string {
	return fmt.Sprintf("MultiMissile %s", Version)
}

func PrintVersion() {
	fmt.Printf(`msl %s
Compiler: %s %s
Copyright (C) 2017 Istyle, Inc.
based on wbt Version %s Copyright (C) 2016 Mercari, Inc.
`,
		Version,
		runtime.Compiler,
		runtime.Version(),
		WbtVersion)
}

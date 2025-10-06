package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type Target struct {
	GOOS   string
	GOARCH string
}

func NewTarget(goos, goarch string) Target {
	if goos == "" {
		goos = runtime.GOOS
	}
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	return Target{
		GOOS:   goos,
		GOARCH: goarch,
	}
}

// LocalTarget returns a target for the local GOOS and GOARCH
func LocalTarget() Target {
	return NewTarget(runtime.GOOS, runtime.GOARCH)
}

func (t Target) IsLocal() bool {
	return t.GOOS == runtime.GOOS && t.GOARCH == runtime.GOARCH
}

func (t Target) Name(p string) string {
	name := filepath.Base(p)
	if !t.IsLocal() {
		name = fmt.Sprintf(
			"%s-%s-%s",
			name,
			t.GOOS,
			t.GOARCH,
		)
	}

	if t.GOOS == "windows" {
		name += ".exe"
	}

	return name
}

func (t Target) Env() map[string]string {
	return map[string]string{
		"GOOS":   t.GOOS,
		"GOARCH": t.GOARCH,
	}
}

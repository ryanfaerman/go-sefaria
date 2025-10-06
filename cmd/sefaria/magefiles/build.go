package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/ctrlc"
	"github.com/caarlos0/log"
	"github.com/charmbracelet/lipgloss"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/ryanfaerman/go-sefaria/cmd/sefaria/magefiles/module"
)

var (
	goexe = "go"
	dirs  = []string{"bin", "tmp"}

	boldStyle = lipgloss.NewStyle().Bold(true)
	codeStyle = lipgloss.NewStyle().Italic(true)
)

// Generate with go generate
func Generate() {
	log.Info("generating code")
	sh.RunV(goexe, "generate", "./...")
}

// Build the sefaria command for local use
func Build() error {
	started := time.Now()
	mg.SerialDeps(ensureDirs, Deps.Tidy, Generate)

	log.Info("building development version(s)")
	log.IncreasePadding()
	defer log.ResetPadding()
	cmd := "sefaria"
	if err := ctrlc.Default.Run(context.Background(), func() error {
		log.IncreasePadding()
		defer log.ResetPadding()

		target := NewTarget(os.Getenv("GOOS"), os.Getenv("GOARCH"))
		name := target.Name(cmd)

		binaryPath := filepath.Join("./bin", name)
		sourcePath := module.Path()

		log.WithField("binary", binaryPath).Infof("building %s", codeStyle.Render(cmd))

		if err := sh.Run(goexe, "build",
			"-o", binaryPath,
			"-buildvcs=false",
			"-tags", "osusergo,netgo",
			"-trimpath",
			sourcePath,
		); err != nil {

			st := log.Styles[log.ErrorLevel]
			log.Warnf("%s %s - %s", st.Render("⚠"), codeStyle.Render(cmd), st.Render("build failed"))
			return fmt.Errorf("failed to build %s: %w", cmd, err)
		}

		return nil
	}); err != nil {
		os.Exit(1)
	}
	// }

	log.Infof("build succeeded after %s", time.Since(started))

	return nil
}

func ensureDirs() error {
	log.Info("preparing output directories")

	log.IncreasePadding()
	defer log.ResetPadding()
	for _, dir := range dirs {
		if !exists("./" + dir) {
			log.WithField("directory", dir).Info("creating")
			if err := os.MkdirAll("./"+dir, 0o755); err != nil {
				return err
			}
		}
	}

	return nil
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

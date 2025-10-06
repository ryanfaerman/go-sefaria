package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caarlos0/ctrlc"
	"github.com/caarlos0/log"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
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

// Build all commands in the cmd directory
func Build() error {
	started := time.Now()
	mg.SerialDeps(ensureDirs, Deps.Tidy, Generate)

	// cmds, err := commands()
	// if err != nil {
	// 	return fmt.Errorf("failed to list products: %w", err)
	// }

	log.Info("building development version(s)")
	log.IncreasePadding()
	defer log.ResetPadding()
	cmd := "sefaria"
	// for _, cmd := range cmds {
	if err := ctrlc.Default.Run(context.Background(), func() error {
		log.IncreasePadding()
		defer log.ResetPadding()

		target := NewTarget(os.Getenv("GOOS"), os.Getenv("GOARCH"))
		name := target.Name(cmd)

		binaryPath := filepath.Join("./bin", name)
		// sourcePath := filepath.Join(module.Path(), "cmd", cmd)
		sourcePath := module.Path()
		spew.Dump(binaryPath, sourcePath)

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

// commands finds all directories within "./cmd" that contain at least one Go file with "package main".
func commands() ([]string, error) {
	var mainDirs []string

	const baseDir = "./cmd"

	err := filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-directories and base directory itself.
		if path == baseDir || !info.IsDir() {
			return nil
		}

		if exists(filepath.Join(path, ".build.skip")) {
			return nil
		}

		// Check if the directory contains at least one Go file with "package main".
		containsMain, err := containsMainPackage(path)
		if err != nil {
			return err
		}
		if containsMain {
			mainDirs = append(mainDirs, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return mainDirs, nil
}

// containsMainPackage checks if a directory contains a Go file with "package main".
func containsMainPackage(dir string) (bool, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		// Check for .go files.
		if strings.HasSuffix(entry.Name(), ".go") {
			filePath := filepath.Join(dir, entry.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return false, err
			}

			if strings.Contains(string(content), "go:build ignore") {
				return false, nil
			}

			// Check if the file contains "package main".
			if strings.Contains(string(content), "package main") {
				return true, nil
			}
		}
	}

	return false, nil
}

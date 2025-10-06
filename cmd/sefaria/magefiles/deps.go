package main

import (
	"strings"

	"github.com/caarlos0/log"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Deps mg.Namespace

// Tidy up current dependencies
func (Deps) Tidy() {
	log.Info("tidying up dependencies")
	sh.RunV("go", "mod", "tidy")
}

// Check dependencies for updates
func (Deps) Check() error {
	log.Infof("Checking for updates... %s", codeStyle.Render("(this might take a while)"))

	// go list -m -u -f '{{if not (or .Indirect .Main)}}{{.Update}}{{end}}' all
	output, err := sh.Output("go", "list", "-u", "-m", "-f", "'{{if not (or .Indirect .Main)}}{{.Version}} {{.Update}}{{end}}'", "all")

	updates := 0

	log.IncreasePadding()
	defer log.ResetPadding()

	for l := range strings.SplitSeq(output, "\n") {
		if l == "''" {
			continue
		}
		if l == "'<nil>'" {
			continue
		}
		l = strings.Trim(l, "'")
		parts := strings.SplitN(l, " ", 3)
		if len(parts) != 3 {
			continue
		}
		updates++
		log.WithField("installed", parts[0]).
			WithField("available", parts[2]).
			Info(codeStyle.Render(parts[1]))

	}
	if updates == 0 {
		log.Info("All dependencies are up to date")
		return nil
	}

	// tbl.Print()
	log.DecreasePadding()

	log.Infof("Found %d updates", updates)
	return err
}

// Update external dependencies
func (Deps) Update() {
	log.Infof("Checking for updates... %s", codeStyle.Render("(this might take a while)"))

	// go list -m -u -f '{{if not (or .Indirect .Main)}}{{.Update}}{{end}}' all
	output, err := sh.Output("go", "list", "-u", "-m", "-f", "'{{if not (or .Indirect .Main)}}{{.Version}} {{.Update}}{{end}}'", "all")
	if err != nil {
		log.Fatal(err.Error())
	}

	log.IncreasePadding()
	defer log.ResetPadding()

	for l := range strings.SplitSeq(output, "\n") {
		if l == "''" {
			continue
		}
		if l == "'<nil>'" {
			continue
		}
		l = strings.Trim(l, "'")
		parts := strings.SplitN(l, " ", 3)
		if len(parts) != 3 {
			continue
		}

		log.Infof("Updating %s from %s to %s", parts[1], parts[0], parts[2])
		out, err := sh.Output("go", "get", "-u", parts[1])
		if err != nil {
			log.Error(out)
			continue
		}
		log.IncreasePadding()
		for o := range strings.SplitSeq(out, "\n") {
			log.Info(o)
		}
		log.DecreasePadding()
	}
}

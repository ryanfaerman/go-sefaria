package main

import (
	"os"

	"github.com/caarlos0/log"
)

// Clean any created directories
func Clean() {
	log.Info("cleaning output directories")
	log.IncreasePadding()

	for _, dir := range dirs {
		if !exists(dir) {
			continue
		}
		log.WithField("directory", dir).Info("removing")
		os.RemoveAll("./" + dir)
	}

	log.DecreasePadding()

	log.Info("done!")
}

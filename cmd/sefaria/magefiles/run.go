package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/caarlos0/log"
	"github.com/magefile/mage/mg"
)

func Run(cmd string) {
	mg.Deps(Build)

	bin := filepath.Join("./bin", cmd)
	log.Infof("%s %s", "running application", codeStyle.Render(bin))

	fmt.Println("---------")

	args := append([]string{cmd}, os.Args[3:]...)
	if err := syscall.Exec(bin, args, os.Environ()); err != nil {
		panic(err.Error())
	}
}

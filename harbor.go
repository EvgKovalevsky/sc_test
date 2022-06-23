package main

import (
	"fmt"
	"os/exec"
)

func pushToHarbor(_filename, _filepath string) error {
	name := _filename + ".tar.gz"
	// command := exec.Command("C:\\Users\\EAKOVALEVS\\bin\\oras.exe",
	command := exec.Command("oras",
		"push",
		fmt.Sprintf("%s/%s/%s:latest", config.HarborHost, config.HarborRepo, _filename),
		fmt.Sprintf("%s:application/vnd.oci.image.layer.v1.tar+gzip", name))
	command.Dir = _filepath
	command.Stdout = logger.Out
	command.Stderr = logger.Out
	command.Run()
	return nil
}

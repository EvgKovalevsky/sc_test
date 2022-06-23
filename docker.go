package main

import (
	"fmt"
	"os"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func startDocker(parser ConfigParser) error {
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return err
	}
	logger.Debug("Docker client created!")

	var containername = parser.Name + "_autoparser"
	cnts, err := client.ListContainers(docker.ListContainersOptions{All: true, Filters: map[string][]string{"name": {containername}}})
	if err != nil {
		return err
	}
	if len(cnts) > 0 && cnts[0].State == "exited" {
		if err := client.RemoveContainer(docker.RemoveContainerOptions{ID: cnts[0].ID, Force: true}); err != nil {
			return err
		}
	}
	logger.Debug("Existing containers founded and deleted!")

	dockerHostConfig := docker.HostConfig{
		VolumeDriver: "bind",
		AutoRemove:   parser.AutoRemove,
		Binds:        []string{parser.ResultHostDir + ":" + parser.ResultDockerDir},
	}

	var command []string = []string{"./app"}
	cf := docker.Config{AttachStdout: true, AttachStderr: true, WorkingDir: "/data", Image: parser.Dockername, Cmd: command}
	opts := docker.CreateContainerOptions{Name: containername, Config: &cf, HostConfig: &dockerHostConfig}

	container, err := client.CreateContainer(opts)
	if err != nil {
		return err
	}
	logger.Debug("Docker container created!")

	err = client.StartContainer(container.ID, nil)
	if err != nil {
		return err
	}
	logger.Debug("Started container with ID:" + container.ID)

	container_logger := &logrus.Logger{
		Out:          os.Stderr,
		Hooks:        map[logrus.Level][]logrus.Hook{},
		Formatter:    &easy.Formatter{TimestampFormat: "2006-01-02 15:04:05", LogFormat: "\t\t%msg%\n"},
		ReportCaller: false,
		Level:        logrus.Level(config.LogLevel),
		ExitFunc: func(int) {
		},
	}

	attachOptsConfig := docker.AttachToContainerOptions{
		Container:   container.ID,
		ErrorStream: container_logger.Writer(),
		Stream:      true,
		Stderr:      true,
	}

	err = client.AttachToContainer(attachOptsConfig)
	if err != nil {
		return err
	}

	code, err := client.WaitContainer(container.ID)
	if err != nil {
		return err
	}
	logger.Debug(fmt.Sprintf("Container finished with code %d", code))

	if code == 0 {

	} else if code == 1 {
		return err
	}
	return nil
}

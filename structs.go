package main

type ConfigParser struct {
	Name            string `mapstructure:"name"`
	Dockername      string `mapstructure:"dockername"`
	AutoRemove      bool   `mapstructure:"auto_remove"`
	ResultHostDir   string `mapstructure:"result_hostdir"`
	ResultDockerDir string `mapstructure:"result_dockerdir"`
}

type ConfigTemplate struct {
	DockerHomeFolder string
	HomeFolder       string
	LogLevel         int
	HarborHost       string
	HarborRepo       string
	Parsers          []ConfigParser
}

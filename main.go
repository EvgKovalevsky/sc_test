package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var logger *logrus.Logger
var config ConfigTemplate

func init() {
	fmt.Println("WELCOME TO GOLANG PARSER MANAGER")

	// Config init
	fmt.Println("Config initialization...")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("File not found!\n")
		} else {
			fmt.Printf("File found but %v!\n", err)
		}
	}
	config = ConfigTemplate{
		DockerHomeFolder: viper.GetString("config.docker_home_folder"),
		HomeFolder:       viper.GetString("config.home_folder"),
		LogLevel:         viper.GetInt("config.log_level"),
		HarborHost:       viper.GetString("config.harbor_host"),
		HarborRepo:       viper.GetString("config.harbor_repo"),
		Parsers:          []ConfigParser{},
	}
	if err := viper.UnmarshalKey("config.parsers", &config.Parsers); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Complete!")

	fmt.Println(config)

	// Logrus init
	logger = &logrus.Logger{
		Out:          os.Stderr,
		Hooks:        map[logrus.Level][]logrus.Hook{},
		Formatter:    &easy.Formatter{TimestampFormat: "2006-01-02 15:04:05", LogFormat: "[%lvl%]\t[%time%] %msg%\n"},
		ReportCaller: false,
		Level:        logrus.Level(config.LogLevel),
		ExitFunc: func(int) {
		},
	}
	logger.Info("Logger initialized!")
}

func main() {
	for {
		updateData()
		time.Sleep(8 * time.Hour)
	}
}

func updateData() {
	for _, parser := range config.Parsers {
		logger.Info(parser.Name + " is parsing...")
		err := startDocker(parser)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"Action":  "Start Container",
				"Message": err,
			}).Error("Docker")
			continue
		}

		files, err := WalkMatch(filepath.Join([]string{config.DockerHomeFolder, "data", parser.Name}...), "*.json")
		if err != nil {
			logger.WithFields(logrus.Fields{
				"Action":  "Walk Match",
				"Message": err,
			}).Error("Files")
		}

		err = createArchive(files, parser.Name+".tar.gz")
		if err != nil {
			logger.WithFields(logrus.Fields{
				"Action":  "Create Archive",
				"Message": err,
			}).Error("Archive")
		}

		logger.Info("Archive created successfully!")
		logger.Info("Push in Harbor")

		err = pushToHarbor(parser.Name, config.DockerHomeFolder)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"Action":  "Push to Harbor",
				"Message": err,
			}).Error("Harbor")
		}

		logger.Info("Pushed")
		fmt.Println()
	}
}

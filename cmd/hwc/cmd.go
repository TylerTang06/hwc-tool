package main

import (
	"os"

	"github.com/TylerTang06/hwc-tool/commands"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"github.com/urfave/cli"
)

func init() {
	initLogger()
}

func main() {
	app := cli.NewApp()
	app.Name = "huaweicloud tool"
	app.Author = "tylertang"
	app.Usage = "Get huaweicloud public resource infos, for example images"
	app.Version = "v1.0"
	app.Commands = commands.Commands
	if err := app.Run(os.Args); err != nil {
		logrus.Error(err)
	}
}

func initLogger() {
	fileName := "logs/hwc_tool.log"
	logrus.SetReportCaller(true)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logrus.SetOutput(file)
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     7,
		Level:      logrus.WarnLevel,
		Formatter: &logrus.TextFormatter{
			DisableTimestamp: false,
			DisableColors:    false,
		},
	})
	logrus.AddHook(rotateFileHook)
}

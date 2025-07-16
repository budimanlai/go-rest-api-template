package main

import (
	"fmt"

	application "go-rest-api-template/internal/application"

	gocli "github.com/budimanlai/go-cli"
)

func main() {
	cli := gocli.NewCliWithConfig(gocli.CliOptions{
		AppName:    "Rest API Service Template",
		Version:    "1.0.0",
		ConfigFile: []string{"configs/config.json"},
	})

	cli.StartService("run", "start", application.RestApi)
	cli.StopService("stop")

	e := cli.Run()
	if e != nil {
		fmt.Println(e.Error())
	}
}

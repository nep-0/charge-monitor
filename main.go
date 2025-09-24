package main

import (
	"charge-monitor/app"
	"charge-monitor/config"
)

func main() {
	config, err := config.ConfigFromFile()
	if err != nil {
		panic(err)
	}
	a := app.NewApp(config.Outlets, config.PollingInterval, config.HTTPAddress)
	a.ServeHTTP()
}

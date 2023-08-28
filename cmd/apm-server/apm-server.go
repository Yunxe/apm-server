package main

import (
	apm_server "APM-server/internal/apm-server"
	"fmt"
	"os"
)

func main() {
	fmt.Println("---------")

	command := apm_server.NewApmServerCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
	fmt.Println("---------")

}

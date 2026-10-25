package main

import (
	"flag"
	"fmt"

	"github.com/LabSketch/Tessera/pkg/config"
)

func main() {
	configDir := flag.String("configDir", "config", "Caminho para pasta de configuração")
	flag.Parse()

	if settings, err := config.Load(*configDir); err == nil {
		message := settings.Application.Message

		fmt.Println(message)
	}
}

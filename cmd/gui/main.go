package main

import (
	"vado/internal/gui"
	"vado/internal/util"
	"vado/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	log, _ := logger.Init()
	defer logger.Sync()

	log.Info("Starting GUI mode.", zap.String("mode", util.GetModeValue()))
	gui.ShowMainApp()
}

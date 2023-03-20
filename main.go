package main

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/goscript-example/lights"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var modules map[string]interface{}
	modules = map[string]interface{}{
		"lights": &lights.Lights{},
	}

	config, err := goscript.ParseConfig("config.yml", modules)
	if err != nil {
		panic(err)
	}

	gs, err := goscript.New(config, setupLogging())
	if err != nil {
		panic(err)
	}

	gs.AddTriggers(lights.Triggers(gs)...)

	err = gs.Connect()
	if err != nil {
		panic(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	gs.GetLogger().Info("Connected to the server")
	<-done

	gs.Close()
}

func setupLogging() logr.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerologr.NameFieldName = "logger"
	zerologr.NameSeparator = "/"
	zerologr.SetMaxV(1)

	zl := zerolog.New(os.Stderr)
	zl = zl.With().Caller().Timestamp().Logger()
	return zerologr.New(&zl)
}

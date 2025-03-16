package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/healthy-heroes/neskuchka/backend/app/cmd"
)

type Opts struct {
	ServerCmd cmd.ServerCommand `command:"server" description:"Start the server"`
}

func main() {
	fmt.Println("Starting application...")

	var opts Opts
	p := flags.NewParser(&opts, flags.Default)

	p.CommandHandler = func(command flags.Commander, args []string) error {
		setupLog()

		// Each command implements the CommonOptionsCommander interface
		c := command.(cmd.CommonOptionsCommander)

		err := c.Execute(args)
		if err != nil {
			log.Error().Err(err).Msg("Command execution failed")
		}

		return err
	}

	if _, err := p.Parse(); err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	fmt.Println("Application finished.")
}

func setupLog() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	log.Logger = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Logger()

	// todo debug mode
	log.Logger = log.Logger.Level(zerolog.DebugLevel)
}

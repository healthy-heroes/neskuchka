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

// Opts contains base options and commands
type Opts struct {
	ServerCmd cmd.ServerCommand `command:"server" description:"Start the server"`

	Debug bool `long:"debug" env:"DEBUG" description:"Enable debug mode"`
}

// revision is set during build
var revision = "unknown"

func main() {
	fmt.Printf("Starting application (revision: %s)\n", revision)

	var opts Opts
	p := flags.NewParser(&opts, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		setupLog(&opts)

		// Each command implements the CommonOptionsCommander interface
		c := command.(cmd.CommonOptionsCommander)
		c.SetCommon(&cmd.CommonOptions{
			Revision: revision,
		})

		err := c.Execute(args)
		if err != nil {
			log.Error().Err(err).Msg("Command execution failed")
		}

		return err
	}

	// Parsing command line arguments and handling errors
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

// setupLog sets up the logger
func setupLog(opts *Opts) {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	log.Logger = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Logger()

	if opts.Debug {
		log.Logger = log.Logger.Level(zerolog.DebugLevel)
	}
}

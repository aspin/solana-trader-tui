package main

import (
	"fmt"
	"github.com/aspin/solana-trader-tui/flags"
	applog "github.com/aspin/solana-trader-tui/log"
	"github.com/aspin/solana-trader-tui/program"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "solana-trader-terminal-ui",
		Usage: "Terminal UI application for interacting with bloXroute Labs's Solana Trader API",
		Flags: []cli.Flag{
			flags.LogFile,
			flags.ConfigFile,
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error while running application: %v",
			err)
	}
}

func run(c *cli.Context) error {
	logConfig := applog.NewConfigFromCLI(c)
	logFile, err := applog.Init(logConfig)
	if err != nil {
		return fmt.Errorf("could not initialize logger: %w", err)
	}
	defer func(logFile *os.File) {
		_ = logFile.Close()
	}(logFile)

	p := program.New()
	_, err = p.Run()
	return err
}

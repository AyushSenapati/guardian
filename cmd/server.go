package cmd

import (
	"context"
	"html"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/AyushSenapati/guardian/lib/server"

	"github.com/AyushSenapati/guardian/config"
	"github.com/spf13/cobra"
)

var (
	configFile         string
	svcDefinitionFname string
	globalConfig       *config.Specification
)

func startServerCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts guardian",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startGuardian(ctx)
		},
	}

	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "guard.json", "Config file")
	cmd.PersistentFlags().StringVarP(
		&svcDefinitionFname, "svcdef", "d", "definitions.json", "service definition file",
	)

	return cmd
}

// it starts the gateway server
func startGuardian(ctx context.Context) error {
	globalConfig, err := config.Load(configFile)
	log.Printf("%+v", globalConfig)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	// Create Guardian server instance with global config
	srv := server.NewServerWithConfig(globalConfig)

	// Create a context which should be canceled
	// on interrupt before starting the server
	ctx = contextWithInterruptSignal(ctx)

	srv.Start(ctx, svcDefinitionFname)
	log.Println(
		"Gaurdian >> now sit back, I am up",
		html.UnescapeString("&#"+strconv.Itoa(128526)+";"),
	)

	srv.Wait(ctx) // waits for stop signal

	log.Println("Guardian >> Bubyee", html.UnescapeString("&#"+strconv.Itoa(9995)+";"))
	return nil
}

// It launches a go routine which listens for Interrupts.
// Once intterupt is received that go routine will cancel the context
func contextWithInterruptSignal(ctx context.Context) context.Context {
	newCtx, cancel := context.WithCancel(ctx)

	// This channel will carry interrupt signal
	signalChan := make(chan os.Signal, 1)

	// Listen and relay the interrupt signal
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		select {
		case <-signalChan:
			cancel()
			close(signalChan)
		}
	}()

	return newCtx
}

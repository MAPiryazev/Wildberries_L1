package main

import (
	"fmt"
	"log"
	"os"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/config"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/cut"
)

func main() {
	flags, err := parseCLI()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if err := config.Init(flags.EnvFile); err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	logger := config.GetLogger()

	fields, err := cut.ParseFields(flags.Fields)
	if err != nil {
		logger.Error().Err(err).Msg("invalid fields")
		os.Exit(1)
	}

	processor, err := cut.NewProcessor(flags.Delimiter, fields, flags.SuppressNoDelim)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create processor")
		os.Exit(1)
	}

	switch flags.Mode {
	case "local":
		if err := processor.ProcessReader(os.Stdin, os.Stdout); err != nil {
			logger.Error().Err(err).Msg("processing error")
			os.Exit(1)
		}
	case "worker":
		log.Fatal("worker mode not yet implemented")
	case "coordinator":
		log.Fatal("coordinator mode not yet implemented")
	default:
		fmt.Fprintf(os.Stderr, "unknown mode: %s\n", flags.Mode)
		os.Exit(1)
	}
}

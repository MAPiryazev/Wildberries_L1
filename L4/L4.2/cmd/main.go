package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/broker"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/config"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/coordinator"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/cut"
	"github.com/MAPiryazev/Wildberries_L1/tree/main/L4/L4.2/internal/worker"
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
	cfg := config.Get()
	if cfg == nil {
		fmt.Fprintln(os.Stderr, "config is not initialized")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	switch flags.Mode {
	case "local":
		processor, _, err := buildProcessor(flags.Delimiter, flags.Fields, flags.SuppressNoDelim)
		if err != nil {
			logger.Error().Err(err).Msg("failed to init local processor")
			os.Exit(1)
		}

		if err := processor.ProcessReader(os.Stdin, os.Stdout); err != nil {
			logger.Error().Err(err).Msg("processing error")
			os.Exit(1)
		}

	case "worker":
		processor, _, err := buildProcessor(cfg.Delimiter, cfg.Fields, cfg.SuppressNoDelim)
		if err != nil {
			logger.Error().Err(err).Msg("failed to init worker processor from config")
			os.Exit(1)
		}

		b, err := broker.NewBroker(cfg.RabbitMQURL, logger)
		if err != nil {
			logger.Error().Err(err).Msg("failed to create broker")
			os.Exit(1)
		}
		defer func() { _ = b.Close() }()

		w := worker.NewWorker(
			cfg.WorkerID,
			b,
			processor,
			cfg.WorkerThreads,
			logger,
		)

		if err := w.Start(ctx); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			logger.Error().Err(err).Msg("worker error")
			os.Exit(1)
		}

	case "coordinator":
		processor, fields, err := buildProcessor(flags.Delimiter, flags.Fields, flags.SuppressNoDelim)
		if err != nil {
			logger.Error().Err(err).Msg("failed to init coordinator processor")
			os.Exit(1)
		}

		b, err := broker.NewBroker(cfg.RabbitMQURL, logger)
		if err != nil {
			logger.Error().Err(err).Msg("failed to create broker")
			os.Exit(1)
		}
		defer func() { _ = b.Close() }()

		coord := coordinator.NewCoordinator(
			processor,
			cfg.QuorumSize,
			b,
			logger,
		)

		if err := coord.ProcessWithQuorum(
			ctx,
			os.Stdin,
			os.Stdout,
			flags.Delimiter,
			fields,
			flags.SuppressNoDelim,
			cfg.ChunkSize,
			cfg.TimeoutQuorum,
		); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return
			}
			logger.Error().Err(err).Msg("coordinator error")
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "unknown mode: %s\n", flags.Mode)
		os.Exit(1)
	}
}

func buildProcessor(delim string, fieldsStr string, suppress bool) (cut.Processor, []int, error) {
	fields, err := cut.ParseFields(fieldsStr)
	if err != nil {
		return nil, nil, err
	}

	p, err := cut.NewProcessor(delim, fields, suppress)
	if err != nil {
		return nil, nil, err
	}

	return p, fields, nil
}

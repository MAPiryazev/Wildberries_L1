package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type CLIFlags struct {
	Delimiter       string
	Fields          string
	SuppressNoDelim bool
	Mode            string
	EnvFile         string
}

func parseCLI() (*CLIFlags, error) {
	fs := pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	fs.String("d", ",", "field delimiter")
	fs.String("f", "", "fields to cut (e.g., 1,3,5-7)")
	fs.BoolP("s", "s", false, "suppress lines without delimiter")
	fs.String("env", ".env", "path to .env file")
	fs.String("mode", "local", "run mode: local, worker, coordinator")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	delim, _ := fs.GetString("d")
	fields, _ := fs.GetString("f")
	suppress, _ := fs.GetBool("s")
	envFile, _ := fs.GetString("env")
	mode, _ := fs.GetString("mode")

	if fields == "" {
		return nil, fmt.Errorf("flag -f (fields) is required")
	}

	return &CLIFlags{
		Delimiter:       delim,
		Fields:          fields,
		SuppressNoDelim: suppress,
		Mode:            mode,
		EnvFile:         envFile,
	}, nil
}

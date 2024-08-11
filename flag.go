package main

import (
	"flag"
	"os"
)

var (
	BreakdownFlag bool
	HelpFlag      bool
	SetupFlag     bool
	VersionFlag   bool
	LogLevel      string
)

func PrintUsage() {
	flag.PrintDefaults()
	os.Exit(1)
}

func initFlags() {
	flag.BoolVar(&BreakdownFlag, "b", false, "breaks down daemon")
	flag.BoolVar(&BreakdownFlag, "breakdown", false, "breaks down daemon")

	flag.BoolVar(&HelpFlag, "h", false, "prints this helpful message")
	flag.BoolVar(&HelpFlag, "help", false, "prints this helpful message")

	flag.BoolVar(&SetupFlag, "s", false, "sets up daemon")
	flag.BoolVar(&SetupFlag, "setup", false, "sets up daemon")

	flag.BoolVar(&VersionFlag, "v", false, "prints version")
	flag.BoolVar(&VersionFlag, "version", false, "prints version")

	flag.StringVar(&LogLevel, "loglevel", "debug", "set log level")
	flag.StringVar(&LogLevel, "l", "debug", "set log level")

	flag.Parse()
}

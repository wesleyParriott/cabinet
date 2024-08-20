package main

import (
	"flag"
	"os"
)

// FIXME: Globals maybe need to get passsed to functions or setup in a state struct or maybe is okay and who cares
var (
	BreakdownFlag bool
	HelpFlag      bool
	PublicFlag    bool
	SetupFlag     bool
	VersionFlag   bool

	LogLevel string
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

	flag.BoolVar(&PublicFlag, "p", false, "runs as public. Meaning no passcode entry when trying to hit the site")
	flag.BoolVar(&PublicFlag, "public", false, "runs as public. Meaning no passcode entry when trying to hit the site")

	flag.BoolVar(&SetupFlag, "s", false, "sets up daemon")
	flag.BoolVar(&SetupFlag, "setup", false, "sets up daemon")

	flag.BoolVar(&VersionFlag, "v", false, "prints version")
	flag.BoolVar(&VersionFlag, "version", false, "prints version")

	flag.StringVar(&LogLevel, "loglevel", "debug", "set log level")
	flag.StringVar(&LogLevel, "l", "debug", "set log level")

	flag.Parse()
}

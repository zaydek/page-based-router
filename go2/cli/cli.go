package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/zaydek/retro/color"
	"github.com/zaydek/retro/loggers"
)

func version() {
	fmt.Println("TODO")
}

func usage() {
	fmt.Println(usageStr)
}

func parseCreateCommandFlags(args []string) *CreateCommandFlags {
	cmd := flag.NewFlagSet("create", flag.ContinueOnError)
	cmd.SetOutput(ioutil.Discard)

	flags := &CreateCommandFlags{}
	cmd.StringVar(&flags.Language, "language", "js", "")
	if err := cmd.Parse(args); err != nil || len(cmd.Args()) > 0 {
		loggers.Stderr.Println(color.Boldf("'retro %s'", strings.Join(os.Args[1:], " ")) + " uses unknown flags and or arguments. " +
			"Try " + color.Bold("'retro help'") + " for help.")
		os.Exit(2)
	}
	if flags.Language != "js" && flags.Language != "ts" {
		loggers.Stderr.Println(color.Bold("'--language'") + " must be " + color.Bold("'js'") + " for JavaScript or " + color.Bold("'ts'") + " for TypeScript.\n\n" +
			"- Try " + color.Bold("'retro create --language=js ...'") + " for JavaScript\n" +
			"- Try " + color.Bold("'retro create --language=ts ...'") + " for TypeScript")
		os.Exit(2)
	}
	return flags
}

func parseWatchCommandFlags(args []string) *WatchCommandFlags {
	cmd := flag.NewFlagSet("watch", flag.ContinueOnError)
	cmd.SetOutput(ioutil.Discard)

	flags := &WatchCommandFlags{}
	cmd.DurationVar(&flags.Poll, "poll", 250*time.Millisecond, "")
	cmd.IntVar(&flags.Port, "port", 8000, "")
	if err := cmd.Parse(args); err != nil || len(cmd.Args()) > 0 {
		loggers.Stderr.Println(color.Boldf("'retro %s'", strings.Join(os.Args[1:], " ")) + " uses unknown flags and or arguments. " +
			"Try " + color.Bold("'retro help'") + " for help.")
		os.Exit(2)
	}
	if flags.Poll < (100*time.Millisecond) && flags.Poll >= (10*time.Second) {
		loggers.Stderr.Println(color.Bold("'--poll'") + " must be between " + color.Bold("'100ms'") + " and " + color.Bold("'10s'") + ".")
		os.Exit(2)
	} else if (flags.Port < 3e3 || flags.Port >= 4e3) && (flags.Port < 5e3 || flags.Port >= 6e3) && (flags.Port < 8e3 || flags.Port >= 9e3) {
		loggers.Stderr.Println(color.Bold("'--port'") + " must be be " + color.Bold("'3XXX'") + " or " + color.Bold("'5XXX'") + " or " + color.Bold("'8XXX'") + ".")
		os.Exit(2)
	}
	for _, each := range cmd.Args() {
		if _, err := os.Stat(each); os.IsNotExist(err) {
			loggers.Stderr.Println("Failed to stat file or directory " + color.Boldf("'%s'", each) + ".")
			os.Exit(2)
		}
	}
	flags.Directories = cmd.Args()
	return flags
}

func parseBuildCommandFlags(args []string) *BuildCommandFlags {
	cmd := flag.NewFlagSet("build", flag.ContinueOnError)
	cmd.SetOutput(ioutil.Discard)

	flags := &BuildCommandFlags{}
	cmd.BoolVar(&flags.Cached, "cached", false, "")
	if err := cmd.Parse(args); err != nil || len(cmd.Args()) > 0 {
		loggers.Stderr.Println(color.Boldf("'retro %s'", strings.Join(os.Args[1:], " ")) + " uses unknown flags and or arguments. " +
			"Try " + color.Bold("'retro help'") + " for help.")
		os.Exit(2)
	}
	return flags
}

func parseServeCommandFlags(args []string) *ServeCommandFlags {
	cmd := flag.NewFlagSet("serve", flag.ContinueOnError)
	cmd.SetOutput(ioutil.Discard)

	flags := &ServeCommandFlags{}
	cmd.IntVar(&flags.Port, "port", 8000, "")
	if err := cmd.Parse(args); err != nil || len(cmd.Args()) > 0 {
		loggers.Stderr.Println(color.Boldf("'retro %s'", strings.Join(os.Args[1:], " ")) + " uses unknown flags and or arguments. " +
			"Try " + color.Bold("'retro help'") + " for help.")
		os.Exit(2)
	}
	if (flags.Port < 3e3 || flags.Port >= 4e3) && (flags.Port < 5e3 || flags.Port >= 6e3) && (flags.Port < 8e3 || flags.Port >= 9e3) {
		loggers.Stderr.Println(color.Bold("'--port'") + " must be be " + color.Bold("'3XXX'") + " or " + color.Bold("'5XXX'") + " or " + color.Bold("'8XXX'") + ".")
		os.Exit(2)
	}
	return flags
}

func ParseCLIArguments() Commands {
	var cmds Commands

	// $ retro
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}

	switch os.Args[1] {

	// $ retro usage
	case "usage":
		fallthrough
	case "--usage":
		usage()
		os.Exit(0)

	// $ retro help
	case "help":
		fallthrough
	case "--help":
		usage()
		os.Exit(0)

	// $ retro version
	case "version":
		fallthrough
	case "--version":
		fallthrough
	case "-v":
		version()
		os.Exit(0)

	// $ retro init
	case "init":
		// var dirname string
		// if len(os.Args) < 3 {
		// 	stderr.Println("It looks like you’re trying to run 'retro init' in the current directory. " +
		// 		"In that case, use '.' explicitly:\n\n" +
		// 		"- retro init .\n\n" +
		// 		"Or\n\n" +
		// 		"- retro init retro-app")
		// 	os.Exit(2)
		// }
		// dirname = os.Args[2]
		cmds.CreateCommand = parseCreateCommandFlags(os.Args[2:])

	// $ retro watch
	case "watch":
		cmds.WatchCommand = parseWatchCommandFlags(os.Args[2:])

	// $ retro build
	case "build":
		cmds.BuildCommand = parseBuildCommandFlags(os.Args[2:])

	// $ retro serve
	case "serve":
		cmds.ServeCommand = parseServeCommandFlags(os.Args[2:])

	default:
		usage()
		os.Exit(2)
	}
	return cmds
}

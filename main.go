package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func Fatal(err string) {
	fmt.Println(err)
	os.Exit(1)
}

func Help(confdir string, models map[string]Model) {
	var fmtModels strings.Builder
	for abbr, model := range models {
		fmt.Fprintf(&fmtModels, "  %-7s-- %s\n", abbr, model.Name)
	}
	fmt.Printf(
		`use LLM directly from terminal.

Usage: minall [command] [options] [...message]

Commands:
  chat     -- start a new chat session
  pipe     -- read from stdin
  trans    -- translate stdin, require translator model
  help     -- print this message

Use "minall [command] -h" to see options of each command.

if command is not matched, all
arguments will be parsed as a message.

Config file is stored in:
  %s

Models defined in config file:
%s
`, confdir, &fmtModels)
}

func readStdin() string {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		Fatal(err.Error())
	}
	return string(input)
}

func main() {
	confdir := GetConfDir()
	config := Load(confdir)
	models := GetModel(config)

	// define flagset for chat and pipe
	FlagSet0 := flag.NewFlagSet("chat & pipe", flag.ExitOnError)
	modelAbbr := FlagSet0.String("m", config.DefaultModel, "specify which model to use")
	systemMsg := FlagSet0.String("s", config.SystemMsg, "specify system message")

	// define flagset for trans
	FlagSet1 := flag.NewFlagSet("trans", flag.ExitOnError)
	modelAbbrT := FlagSet1.String("m", config.DefaultTranslator, "specify which model to use")
	targetLang := FlagSet1.String("t", "Chinese", "specify target language using its full English name")
	domain := FlagSet1.String("d", "", "describe the domain of the text")

	// define subcommands
	if len(os.Args) < 2 {
		Help(confdir, models)
		os.Exit(0)
	}

	subcommand := os.Args[1]
	switch subcommand {
	case "chat": // start a new chat session
		FlagSet0.Parse(os.Args[2:])
		if !models.IsValidModel(*modelAbbr, []string{"chat", "reasoner"}) {
			Fatal("Invalid Model: " + *modelAbbr)
		}
		chatSession(models[*modelAbbr], *systemMsg)

	case "pipe": // read from stdin
		FlagSet0.Parse(os.Args[2:])
		if !models.IsValidModel(*modelAbbr, []string{"chat", "reasoner"}) {
			Fatal("Invalid Model: " + *modelAbbr)
		}
		Quest(models[*modelAbbr], []Message{
			{"system", *systemMsg},
			{"user", readStdin()},
		})

	case "trans": // translate stdin
		FlagSet1.Parse(os.Args[2:])
		if !models.IsValidModel(*modelAbbrT, []string{"translator"}) {
			Fatal("Invalid Model: " + *modelAbbrT)
		}
		Translate(models[*modelAbbrT], *targetLang, *domain, readStdin())

	case "help":
		Help(confdir, models)

	default: // all arguments are parsed as a message
		FlagSet0.Parse(os.Args[1:])
		if !models.IsValidModel(*modelAbbr, []string{"chat", "reasoner"}) {
			Fatal("Invalid Model: " + *modelAbbr)
		}
		msg := strings.Join(FlagSet0.Args(), " ")
		Quest(models[*modelAbbr], []Message{
			{"system", *systemMsg},
			{"user", msg},
		})
	}
}

package cli

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/chzyer/readline"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("add"),
	readline.PcItem("connect"),
	readline.PcItem("disconnect"),
	readline.PcItem("info"),
	readline.PcItem("invoke"),
	readline.PcItem("link"),
	readline.PcItem("remove"),
	readline.PcItem("set"),
	readline.PcItem("unlink"),
	readline.PcItem("quit"),
	readline.PcItem("notify"),
)

func Run() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:            "> ",
		HistoryFile:       "/tmp/olink.tmp",
		HistorySearchFold: true,
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "quit",
	})
	rl.CaptureExitSignal()
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		words := strings.Split(strings.TrimSpace(line), " ")
		if len(words) == 0 {
			continue
		}
		name := words[0]
		var cmd Command
		for _, c := range commands {
			for _, n := range c.Names {
				if n == name {
					cmd = c
				}
			}
		}
		if cmd.Exec == nil {
			fmt.Printf("unknown command: %s\n", name)
			continue
		} else {
			err := cmd.Exec(words)
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("%s\n", err)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

}

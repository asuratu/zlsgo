package zcli

import (
	"flag"
	"fmt"
	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/zstring"
	"github.com/sohaha/zlsgo/ztype"
	"strings"
)

func init() {
	Log = zlog.New()
	Log.ResetFlags(zlog.BitLevel)
	// flag.CommandLine.SetOutput(ioutil.Discard)
	flag.CommandLine.SetOutput(&errWrite{})
	flag.Usage = func() {
		usage()
	}
}

// Add Registers a Cmd for the provided subCommand name
func Add(name, description string, command Cmd) (cmd *cmdCont) {
	if name == "" {
		Log.Error(getLangs("command_empty"))
		return
	}
	cmd = &cmdCont{
		name:          name,
		desc:          description,
		command:       command,
		requiredFlags: RequiredFlags{},
	}
	cmds[name] = cmd

	return
}

func SetUnknownCommand(fn func(_ string)) {
	unknownCommandFn = fn
}

func usage() {
	showHeadr()
	showFlagsAndRequired := func() {
		if numOfGlobalFlags() > 0 {
			ShowFlags(flag.CommandLine)
			ShowRequired(flag.CommandLine, requiredFlags)
		}
	}
	if len(cmds) == 0 {
		Log.Printf("usage of %s:\n", showText(firstParameter))
		showFlagsAndRequired()
		return
	}
	Log.Printf("usage: %s <command>\n\n", firstParameter)
	Log.Println("  where <command> is one of:")
	for name, cont := range cmds {
		Log.Printf("    "+tipText("%-19s")+" %s\n", name, cont.desc)
	}

	showFlagsAndRequired()
	if !HidePrompt {
		Log.Printf(showText("\nMore Command information, please use: %s <command> --help\n"), firstParameter)
	}
}

func ShowFlags(fg *flag.FlagSet) {
	Log.Printf("\n  optional flags:\n")
	max := 40
	showFlagsHelp()
	flagsItems := zstring.Buffer()
	fg.VisitAll(func(f *flag.Flag) {
		s := zstring.Buffer()
		flagsTitle := f.Name
		output := false
		if flagsTitle == "version" {
			output = true
		}
		name, usage := flag.UnquoteUsage(f)
		// if name == "" {
		// 	name = "bool"
		// }
		sf := "    -%-12s"
		if len(name) > 0 {
			newName := showText("<" + name + ">")
			namePadLen := 12 + len(newName) - len(name)
			flagsTitle += " " + newName
			sf = "    -%-" + ztype.ToString(namePadLen) + "s"
		}
		s.WriteString(warnText(fmt.Sprintf(sf, flagsTitle)))
		if zstring.Len(s.String()) <= max {
			s.WriteString("\t")
		} else {
			s.WriteString("\n    \t")
		}
		s.WriteString(strings.ReplaceAll(usage, "\n", "\n    \t"))
		defValue := ztype.ToString(f.DefValue)
		if defValue != "" && defValue != "0" && defValue != "false" {
			s.WriteString(fmt.Sprintf(" (default %v)", defValue))
		}
		if output {
			Log.Println(s.String())
		} else {
			s.WriteString("\n")
			flagsItems.WriteString(s.String())
		}
	})

	Log.Println(flagsItems.String())
}

// Start Start
func Start(runFunc ...runFunc) {
	if matchingCmd != nil {
		if *flagHelp {
			showSubcommandUsage(flag.CommandLine, matchingCmd)
		} else {
			matchingCmd.command.Run(args)
		}
		return
	}
	isRunFunc := len(runFunc) > 0
	if isRunFunc {
		runFunc[0]()
	} else {
		Help()
	}
}

// Run runnable
func Run(runFunc ...runFunc) (ok bool) {
	isRunFunc := len(runFunc) > 0
	parse(!isRunFunc)
	Start(runFunc...)
	return
}

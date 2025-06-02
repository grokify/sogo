package cobraexample

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/sogo/flag/cobrautil"
	"github.com/spf13/cobra"
)

type ExampleOptions struct {
	PropBool               []bool         `short:"b" long:"bool" description:"Set boolean"`
	PropLong               uint           `long:"long" description:"Long Int"`
	PropCall               func(string)   `short:"c" description:"Callback func" json:"-"`
	PropStringRequired     string         `short:"r" long:"stringrequired" description:"String required" required:"true"`
	PropString             string         `long:"string" choice:"foo" choice:"bar"`
	PropStringValueName    string         `short:"v" long:"valuename" description:"A value-name" value-name:"VALUE"`
	PropPointer            *int           `short:"p" description:"A pointer to an integer"`
	PropStringSlice        []string       `short:"s" description:"A slice of strings"`
	PropStringPointerSlice []*string      `long:"ptrslice" description:"A slice of pointers to string"`
	PropIntMap             map[string]int `long:"intmap" description:"A map from string to int"`
	PropIntSlice           []int          `long:"intslice" default:"1" default:"2" env:"VALUES"  env-delim:","`
}

func (opts *ExampleOptions) RunCobraFunc(cmd *cobra.Command, args []string) {
	if err := opts.RunCobra(cmd, args); err != nil {
		slog.Error("error running cobra command", "errorMessage", err.Error())
		os.Exit(1)
	}
}

func (opts *ExampleOptions) RunCobra(cmd *cobra.Command, args []string) error {
	if cmd == nil {
		return errors.New("cobra.Command cannot be nil")
	} else if err := cobrautil.CobraToGoflagsRuntime(cmd, opts); err != nil {
		return err
	} else {
		opts.Run()
		return nil
	}
}

func (opts *ExampleOptions) Run() {
	fmtutil.MustPrintJSON(opts)
}

func ExampleCommand(cmdName string) (*cobra.Command, error) {
	cmdName = strings.TrimSpace(cmdName)
	if cmdName == "" {
		cmdName = "testtypes"
	}
	opts := &ExampleOptions{}
	var mergeCmd = &cobra.Command{
		Use:   cmdName,
		Short: "Test param types",
		Long:  `Test param types for cobra and go-flags integration.`,
		Run:   opts.RunCobraFunc,
	}

	if err := cobrautil.GoflagsToCobraConfig(mergeCmd, opts); err != nil {
		return nil, err
	} else {
		return mergeCmd, nil
	}
}

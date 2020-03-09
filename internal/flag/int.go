package flag

import "github.com/urfave/cli/v2"

// OptionInt params for build int flag.
type OptionInt func(f *cli.IntFlag)

// IntDefault set default value for IntFlag.
func IntDefault(value int) OptionInt {
	return func(f *cli.IntFlag) {
		f.Value = value
		f.HasBeenSet = true
	}
}

// IntAliases set aliases for IntFlag.
func IntAliases(values ...string) OptionInt {
	return func(f *cli.IntFlag) {
		f.Aliases = values
	}
}

// IntRequired set flag mandatory for IntFlag.
func IntRequired() OptionInt {
	return func(f *cli.IntFlag) {
		f.Required = true
	}
}

// IntEnv set env name for IntEnv.
func IntEnv(values ...string) OptionInt {
	return func(f *cli.IntFlag) {
		f.EnvVars = values
	}
}

// NewIntFlag create new IntFlag by options.
func NewIntFlag(name, usage string, options ...OptionInt) *cli.IntFlag {
	flag := &cli.IntFlag{
		Name:  name,
		Usage: usage,
	}

	for i := range options {
		options[i](flag)
	}

	return flag
}

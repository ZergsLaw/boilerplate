// Package flag need for convenient get object flag from lib cli.
package flag

import (
	"github.com/urfave/cli/v2"
)

// OptionString params for build string flag.
type OptionString func(f *cli.StringFlag)

// StrDefault set default value for StringFlag.
func StrDefault(value string) OptionString {
	return func(f *cli.StringFlag) {
		f.Value = value
		f.HasBeenSet = true
	}
}

// StrAliases set aliases for StringFlag.
func StrAliases(values ...string) OptionString {
	return func(f *cli.StringFlag) {
		f.Aliases = values
	}
}

// StrRequired set flag mandatory for StringFlag.
func StrRequired() OptionString {
	return func(f *cli.StringFlag) {
		f.Required = true
	}
}

// StrEnv set env name for StringFlag.
func StrEnv(values ...string) OptionString {
	return func(f *cli.StringFlag) {
		f.EnvVars = values
	}
}

// NewStrFlag create new StringFlag by options.
func NewStrFlag(name, usage string, options ...OptionString) *cli.StringFlag {
	flag := &cli.StringFlag{
		Name:  name,
		Usage: usage,
	}

	for i := range options {
		options[i](flag)
	}

	return flag
}

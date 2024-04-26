package managedflag

import (
	"github.com/spf13/cobra"
)

type baseFlag struct {
	cmd *cobra.Command
	name string
}

type StrFlag struct {
	baseFlag
	Value *string
}

type Int64Flag struct {
	baseFlag
	Value *int64
}

type BoolFlag struct {
	baseFlag
	Value *bool
}

func (f *baseFlag)IsChanged() bool {
	return f.cmd.Flags().Changed(f.name)
}

func NewStr(cmd *cobra.Command, name string, defaultValue string, usage string) (*StrFlag) {
	var value *string = new(string)
	cmd.Flags().StringVar(value, name, defaultValue, usage)
	return &StrFlag{ baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewStrP(cmd *cobra.Command, name string, shorthand string, defaultValue string, usage string) (*StrFlag) {
	var value *string = new(string)
	cmd.Flags().StringVarP(value, name, shorthand, defaultValue, usage)
	return &StrFlag{ baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewInt64(cmd *cobra.Command, name string, defaultValue int64, usage string) (*Int64Flag) {
	var value *int64 = new(int64)
	cmd.Flags().Int64Var(value, name, defaultValue, usage)
	return &Int64Flag{ baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewInt64P(cmd *cobra.Command, name string, shorthand string, defaultValue int64, usage string) (*Int64Flag) {
	var value *int64 = new(int64)
	cmd.Flags().Int64VarP(value, name, shorthand, defaultValue, usage)
	return &Int64Flag{ baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewBool(cmd *cobra.Command, name string, defaultValue bool, usage string) (*BoolFlag) {
	var value *bool = new(bool)
	cmd.Flags().BoolVar(value, name, defaultValue, usage)
	return &BoolFlag{ baseFlag: baseFlag{cmd, name}, Value: value}
}

func NewBoolP(cmd *cobra.Command, name string, shorthand string, defaultValue bool, usage string) (*BoolFlag) {
	var value *bool = new(bool)
	cmd.Flags().BoolVarP(value, name, shorthand, defaultValue, usage)
	return &BoolFlag{ baseFlag: baseFlag{cmd, name}, Value: value}
}
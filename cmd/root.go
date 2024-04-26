/*
Copyright © 2024 Eugene Shtoka <eshtoka@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	logger zerolog.Logger
	logLevel string
	serviceName = "figoro"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "figoro",
	Version: "0.0.1",
	Short: "List events from multiple Google Calendars, offering customizable filtering.",
	Long: `List events from multiple Google Calendars, offering customizable filtering. 
For example:

figoro --calendar "Work" --start "2023-12-25" --end "2024-01-01"

figoro --mode "json" --start "2023-12-25" --end "2024-01-01".`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initRootCmdConfig)

	// Find home directory.
	home, err := os.UserHomeDir()
	var message, cfgDefault string
	if (err != nil) {
		message = "failed to retrieve user home dir. You will have to set path to config file manually"
	} else {
		message = "path for config file"
		cfgDefault = fmt.Sprintf("%s/.config/figoro/figoro.yaml", home)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgDefault, message)
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "disabled", "log level [debug, info, warn, error, disabled]")

	registerCommands()
}

func registerCommands() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	listCmd.AddCommand(listAccountsCmd)
}

// initConfig reads in config file and ENV variables if set.
func initRootCmdConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)

		// If a config file is found, read it in.
		err := viper.MergeInConfig()
		if err != nil {
			showError(fmt.Sprintf("failed to read config file: %s", cfgFile), err)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		showError(fmt.Sprintf("invalid config for log-level: %s", logLevel), err)
		level = zerolog.Disabled
	}

	logger = zerolog.New(os.Stderr).With().Timestamp().Logger().Level(level)
	logger.Debug().Msgf("reading configuration from: %s\n", viper.ConfigFileUsed())
}

// TODO: fix command documentation
// TODO: add test cases
// TODO: add agenda command
// TODO: add event commands: add, delete, update
// TODO: add flags/filters for list events command
// TODO: support multiple inner calendars 
// TODO: rename `add calendar` to `add account`
// TODO: build CI/CD for the project
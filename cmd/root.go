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

var rootCmd = &cobra.Command{
	Use:   "figoro",
	Version: "0.0.1",
	Short: "View and manage multiple Google Calendars.",
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

// TODO: fix list events documentation
// TODO: add test cases
// TODO: add event commands: add, delete, update
// TODO: build CI/CD for the project
// TODO: add commands to manage whitelist & blacklist of calendars
// TODO: add optional flags to limit events to specific list of calendars
// TODO: extract part of the code to external packages
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
	"slices"

	"github.com/EugeneShtoka/figoro/lib/typedkeyring"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteAccountCmd represents the delete command
var deleteAccountCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := deleteAccountFromConfig(args[0])
		if (err != nil) {
			showError(fmt.Sprintf("failed to delete account '%s'", args[0]), err)
			cmd.Usage()
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteAccountCmd)
}

func deleteAccountFromConfig(calName string) error {
	accounts := []string (viper.GetStringSlice(accountsConfigKey))
	accounts = slices.DeleteFunc(accounts, func(s string) bool { return s == calName })

	viper.Set(accountsConfigKey, accounts)
	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("failed to delete account '%s' to config: %w", calName, err)
	}
	keyring := typedkeyring.New[any](serviceName)
	return keyring.Delete(calName)
}

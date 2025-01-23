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

	"github.com/EugeneShtoka/figoro/lib/gaccount"
	"github.com/EugeneShtoka/figoro/lib/typedkeyring"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteAccountCmd = &cobra.Command{
	Use:   "account [account name]",
	Args:  cobra.ExactArgs(1),
	Short: "Delete account",
	Long: "Delete account. Requires account name to delete",
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

func deleteAccountFromConfig(accName string) error {
	accounts := getAccountsFromConfig()

	predicate := func(acc gaccount.GAccount) bool { return acc.Name == accName }
	if (!slices.ContainsFunc(accounts, predicate)) {
		return fmt.Errorf("account '%s' does not exist in config", accName)	
	}
	accounts = slices.DeleteFunc(accounts, predicate)

	viper.Set(accountsConfigKey, accounts)
	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("failed to delete account '%s' from config: %w", accName, err)
	}
	keyring := typedkeyring.New[any](serviceName)
	return keyring.Delete(accName)
}

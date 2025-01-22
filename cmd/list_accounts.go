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
	"strings"

	"github.com/EugeneShtoka/figoro/lib/gaccount"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List accounts",
	Long: " Displays a list of all Google Calendar accounts that have been authorized and configured to access your calendar data.",
	Run: func(cmd *cobra.Command, args []string) {
		listAccountsFromConfig()
	},
}

func init() {
	listCmd.AddCommand(listAccountsCmd)
}

func listAccountsFromConfig() {
	var accounts []gaccount.GAccount
	err := viper.UnmarshalKey(accountsConfigKey, &accounts)
	if err != nil {
		fmt.Errorf("failed to read accounts from config: %v", err)
	}

	accountsNames := make([]string, len(accounts))
    for i, account := range accounts {
        accountsNames[i] = account.Name
    }
	fmt.Printf("Authorized accounts: %s\n", strings.Join(accountsNames, ", "))
}

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

	"github.com/spf13/cobra"
)

var listAccountsCmd = &cobra.Command{
	Use: "accounts",
	Short: "List accounts",
	Long: "Display a list of all accounts that have been authorized and configured to access your calendar data.",
	Run: func(cmd *cobra.Command, args []string) {
		listAccountsFromConfig()
	},
}

func init() {
	listCmd.AddCommand(listAccountsCmd)
}

func listAccountsFromConfig() {
	accounts := getAccountsFromConfig()
	//accountsNames := xiter.Map(getAccountsIterFromConfig(), func(acc gaccount.GAccount) string { return acc.Name })
	//fmt.Printf("Authorized accounts: %s\n", strings.Join(xiter.ToSlice(accountsNames), ", "))

	for _, account := range accounts {
		fmt.Printf("Account %s\n", account.Name)
		fmt.Printf("Showing info for calendars:\n")
		for	_, name := range account.ResolveCalendars() {
			fmt.Printf("\t%s\n", name)
		}
		fmt.Printf("Available calendars:\n")
		for	_, name := range account.Calendars.All {
			fmt.Printf("\t%s\n", name)
		}
		if (len(account.Calendars.WhiteList) > 0) {
			fmt.Printf("Whitelisted calendars:\n")
			for	_, calendar := range account.Calendars.WhiteList {
				fmt.Printf("\t%s\n", calendar)
			}
		}
		if (len(account.Calendars.BlackList) > 0) {
			fmt.Printf("Blacklisted calendars:\n")
			for	_, calendar := range account.Calendars.BlackList {
				fmt.Printf("\t%s\n", calendar)
			}
		}
	}
}

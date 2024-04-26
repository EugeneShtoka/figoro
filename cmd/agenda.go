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
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/EugeneShtoka/figoro/lib/combaccount"
	"github.com/EugeneShtoka/figoro/lib/eventsfilter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/calendar/v3"
)

// agendaCmd represents the agenda command
var agendaCmd = &cobra.Command{
	Use:   "agenda",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		events, err := agenda()
		if (err != nil) {
			showError("failed to list events", err)
			cmd.Usage()
		}
		fmt.Println(events)
	},
}

func init() {
	rootCmd.AddCommand(agendaCmd)
}


func formatEventsAsAgenda(events []*calendar.Event) (string, error) {
	if len(events) == 0 {
		return "[]", nil
	}
	
	jsonData, err := json.Marshal(events)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func agenda() (string, error) {
		accounts := viper.GetStringSlice(accountsConfigKey)
		ctx := context.Background()

		account, err := combaccount.New(ctx, serviceName, accounts)
		if (err != nil) {
			return "", fmt.Errorf("failed to initialize accounts: %s: %v", accounts, err)
		}

		tFrom := time.Now().Format(time.RFC3339)
		filter := eventsfilter.New().ShowSingle().MinEndTime(tFrom)
		events, err := account.Events(filter)
		if (err != nil) {
			return "", fmt.Errorf("failed to retrieve events for accounts: %s: %v", accounts, err)
		}
		
		eventsString, err := formatEventsAsAgenda(events)
		if (err != nil) {
			return "", fmt.Errorf("failed to convert events to string for accounts: %s: %v", accounts, err)
		}

		return eventsString, nil
}

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
	"github.com/EugeneShtoka/figoro/lib/gaccount"
	"github.com/EugeneShtoka/figoro/lib/managedflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/api/calendar/v3"
)

var (
	minEndTime		*managedflag.StrFlag
	maxStartTime	*managedflag.StrFlag
	eventTypes		*managedflag.StrFlag
	orderBy			*managedflag.StrFlag

	maxResults		*managedflag.Int64Flag

	single			*managedflag.BoolFlag
	deleted			*managedflag.BoolFlag
)

// listEventsCmd represents the events command
var listEventsCmd = &cobra.Command{
	Use:   "events",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		events, err := listEvents()
		if (err != nil) {
			showError("failed to list events", err)
			cmd.Usage() 
		}
		fmt.Println(events)
	},
}

func init() {
	listCmd.AddCommand(listEventsCmd)

	minEndTime = managedflag.NewStr(listEventsCmd, "minEndTime", "", "list events with end times later than (default now)")
	maxStartTime = managedflag.NewStr(listEventsCmd, "maxStartTime", "", "list events with start times earlier than")
	eventTypes = managedflag.NewStr(listEventsCmd, "eventTypes", "", "list events with specified event types")
	orderBy = managedflag.NewStr(listEventsCmd, "orderBy", "", "list events with specified order")

	maxResults = managedflag.NewInt64(listEventsCmd, "maxResults", 0, "max results per account")

	single = managedflag.NewBool(listEventsCmd, "single", false, "list events for all accounts")
	deleted = managedflag.NewBool(listEventsCmd, "deleted", false, "list events for all accounts")
}

func eventsToString(events []*calendar.Event) (string, error) {
	if len(events) == 0 {
		return "[]", nil
	}
	
	jsonData, err := json.Marshal(events)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func listEvents() (string, error) {
		var accounts []*gaccount.GAccount
		err := viper.UnmarshalKey(accountsConfigKey, &accounts)
		if (err != nil) {
			return "", fmt.Errorf("failed to read accounts from config: %v", err)
		}

		ctx := context.Background()

		account, err := combaccount.New(ctx, serviceName, accounts)
		if (err != nil) {
			return "", fmt.Errorf("failed to initialize accounts: %v: %v", accounts, err)
		} 

		minEndTimeValue := time.Now().Format(time.RFC3339)
		if (minEndTime.IsChanged()) {
			minEndTimeValue = *minEndTime.Value
		}
		filter := eventsfilter.New().MinEndTime(minEndTimeValue)

		if (maxStartTime.IsChanged()) {
			filter = filter.MaxStartTime(*maxStartTime.Value)
		}

		if (eventTypes.IsChanged()) {
			filter = filter.EventTypes(*eventTypes.Value)
		}

		if (orderBy.IsChanged()) {
			filter = filter.OrderBy(*orderBy.Value)
		}

		if (maxResults.IsChanged()) {
			filter = filter.MaxResults(*maxResults.Value)
		}

		if (single.IsChanged() && *single.Value) {
			filter = filter.ShowSingle()
		}

		if (deleted.IsChanged() && *deleted.Value) {
			filter = filter.ShowDeleted()
		}

		events, err := account.Events(filter)
		if (err != nil) {
			return "", fmt.Errorf("failed to retrieve events for accounts: %v: %v", accounts, err)
		}

		//ToDo: reorder combined list of events + limit to maxResults
		
		eventsString, err := eventsToString(events)
		if (err != nil) {
			return "", fmt.Errorf("failed to convert events to string for accounts: %v: %v", accounts, err)
		}

		return eventsString, nil
}

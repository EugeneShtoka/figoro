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
package combaccount

import (
	"context"

	"github.com/EugeneShtoka/figoro/lib/eventsfilter"
	"github.com/EugeneShtoka/figoro/lib/gaccount"
	"google.golang.org/api/calendar/v3"
)

type CombinedAccount struct {
	accounts []*gaccount.GAccount
}

func New(ctx context.Context, serviceName string, calNames []string) (*CombinedAccount, error) {
	var accounts []*gaccount.GAccount
	for _, calName := range calNames {
		gAcc, err := gaccount.New(ctx, serviceName, calName)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, gAcc)
	}

	return &CombinedAccount{ accounts }, nil
}

func (ca *CombinedAccount) Events(filter *eventsfilter.EventsFilter) ([]*calendar.Event, error) {
	var combinedEvents []*calendar.Event
	for _, gAcc := range ca.accounts {
		events,  err := gAcc.Events(filter)
		if (err == nil) {
			combinedEvents = append(combinedEvents, events...)
		} else {
			return nil, err
		}
	}
	return combinedEvents, nil	
}
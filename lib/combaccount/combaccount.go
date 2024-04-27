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
	"fmt"
	"sort"
	"time"

	"github.com/EugeneShtoka/figoro/lib/eventsfilter"
	"github.com/EugeneShtoka/figoro/lib/gaccount"
	"google.golang.org/api/calendar/v3"
)

type CombinedAccount struct {
	accounts []*gaccount.GAccount
}

func New(serviceName string, accounts []*gaccount.GAccount) (*CombinedAccount, error) {
	for _, account := range accounts {
		err := account.Init(serviceName)
		if err != nil {
			return nil, err
		}
	}

	return &CombinedAccount{ accounts }, nil
}

func sortEvents(events []*calendar.Event) {
	sort.Slice(events, func(i, j int) bool {
		if (events[i].Start.Date == "") {
			if (events[j].Start.Date == "") {
				return events[i].Start.DateTime < events[j].Start.DateTime
			} else {
				return events[i].Start.DateTime < events[j].Start.Date
			}
		} else if (events[j].Start.Date == "") {
			return events[i].Start.Date < events[j].Start.DateTime 
		}
		return events[i].Start.Date < events[j].Start.Date
	})
}

func (ca *CombinedAccount) Events(filter *eventsfilter.EventsFilter) ([]*calendar.Event, error) {
	start := time.Now()
	_, cancel := context.WithCancel(context.Background())
    defer cancel()
	eventsChannel := make(chan []*calendar.Event)
	errorChanel := make(chan error)

	calCount := 0
	for _, gAcc := range ca.accounts {
		calendars := gAcc.ResolveCalendars()
		calCount += len(calendars)
		for _, calendar := range calendars {
			go func(gAcc *gaccount.GAccount, calendar string) {
				events, err := gAcc.Events(calendar, filter)
				if err != nil {
					errorChanel <- err
					cancel()
				}
				eventsChannel <- events
			} (gAcc, calendar) 
		}
	}

	var combinedEvents []*calendar.Event
	for i := 0; i < calCount; i++ {
		select {
			case events := <-eventsChannel:
				combinedEvents = append(combinedEvents, events...)
			case err := <-errorChanel:
				return nil, fmt.Errorf("failed to get events: %w", err)
		}
	}

	if filter.IsOrderedByStartTime() {
		sortEvents(combinedEvents)
	}

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
	
	return combinedEvents, nil	
}
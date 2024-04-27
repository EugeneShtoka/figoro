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

	"github.com/EugeneShtoka/figoro/lib/concurrentresult"
	"github.com/EugeneShtoka/figoro/lib/eventsfilter"
	"github.com/EugeneShtoka/figoro/lib/gaccount"
	"github.com/EugeneShtoka/figoro/lib/sliceutils"
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

func sortEventsByStartTime(events []*calendar.Event) {
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

func sortEventsByUpdated(events []*calendar.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Updated < events[j].Updated
	})
}

func reapplyFiltersOnCombinedEvents(events []*calendar.Event, filter *eventsfilter.EventsFilter) []*calendar.Event {
	if filter.IsOrderedByStartTime() {
		sortEventsByStartTime(events)
	}

	if filter.IsOrderedByUpdated() {
		sortEventsByUpdated(events)
	}

	maxResults := filter.GetMaxResults()
	if (maxResults != nil) {
		events = events[:*maxResults]
	}

	return events
}

func getEvents(gAcc *gaccount.GAccount, calendar string,  filter *eventsfilter.EventsFilter, concurrentResult *concurrentresult.ConcurrentResult[[]*calendar.Event]) {
	events, err := gAcc.Events(calendar, filter)
	if err != nil {
		concurrentResult.SendError(err)
		concurrentResult.Cancel()
	}
	concurrentResult.SendResult(events)
}

func (ca *CombinedAccount) Events(filter *eventsfilter.EventsFilter) ([]*calendar.Event, error) {
	concurrentResult := concurrentresult.New[[]*calendar.Event](context.Background())
	defer concurrentResult.Cancel()

	calCount := 0
	for _, gAcc := range ca.accounts {
		calendars := gAcc.ResolveCalendars()
		calCount += len(calendars)
		for _, calendar := range calendars {
			go getEvents(gAcc, calendar, filter, concurrentResult)
		}
	}

	events2DArr, err := concurrentResult.Results(calCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	
	combinedEvents := sliceutils.FlattenSlice(events2DArr)
	filteredEvents := reapplyFiltersOnCombinedEvents(combinedEvents, filter)

	return filteredEvents, nil	
}
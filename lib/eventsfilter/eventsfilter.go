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
package eventsfilter

import (
	"google.golang.org/api/calendar/v3"
)

type EventsFilter struct {
	minEndTime		*string
	maxStartTime	*string
	maxResults		*int64
	eventTypes		*string
	orderBy			*string
	single			bool
	deleted			bool
}

func New() *EventsFilter {
	return &EventsFilter{
		minEndTime: nil,
		maxStartTime: nil,
		maxResults: nil,
		eventTypes: nil,
		orderBy: nil,
		single: false,
		deleted: false,
	}
}

func (ef *EventsFilter) MinEndTime (time string) *EventsFilter {
	val := time
	ef.minEndTime = &val
	return ef
}

func (ef *EventsFilter) MaxStartTime (time string) *EventsFilter {
	val := time
	ef.maxStartTime = &val
	return ef
}

func (ef *EventsFilter) MaxResults (results int64) *EventsFilter {
	val := results
	ef.maxResults = &val
	return ef
}

func (ef *EventsFilter) EventTypes (types string) *EventsFilter {
	val := types
	ef.eventTypes = &val
	return ef
}

func (ef *EventsFilter) OrderBy (order string) *EventsFilter {
	val := order
	ef.orderBy = &val
	return ef
}

func (ef *EventsFilter) ShowSingle () *EventsFilter {
	ef.single = true
	return ef
}

func (ef *EventsFilter) ShowDeleted () *EventsFilter {
	ef.deleted = true
	return ef
}

func (ef *EventsFilter) IsOrderedByStartTime () bool {
	return (*ef.orderBy == "startTime")
}

func (ef *EventsFilter) Apply (listCall *calendar.EventsListCall) *calendar.EventsListCall {
	listCall = listCall.ShowDeleted(ef.deleted).SingleEvents(ef.single);

	if (ef.minEndTime != nil) {
		listCall = listCall.TimeMin(*ef.minEndTime)
	}

	if (ef.maxStartTime != nil) {
		listCall = listCall.TimeMax(*ef.maxStartTime)
	}

	if (ef.maxResults != nil) {
		listCall = listCall.MaxResults(*ef.maxResults)
	}

	if (ef.eventTypes != nil) {
		listCall = listCall.EventTypes(*ef.eventTypes)
	}

	if (ef.orderBy != nil) {
		listCall = listCall.OrderBy(*ef.orderBy)
	}

	return listCall
}
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
package gaccount

import (
	"context"
	"fmt"

	"github.com/EugeneShtoka/figoro/lib/eventsfilter"
	"github.com/EugeneShtoka/figoro/lib/gaseed"
	"github.com/EugeneShtoka/figoro/lib/typedkeyring"
	set "github.com/deckarep/golang-set/v2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GCalendars struct {
	All 		[]string
	WhiteList	[]string
	BlackList	[]string
}

type GAccount struct {
	Name 			string
	Calendars 		GCalendars
	Service 		*calendar.Service
}

func New(serviceName string, accountName string) (*GAccount, error) {
	service, err := getService(serviceName, accountName)
	if (err != nil) {
		return nil, err
	}

	calendars, err := getCalendars(service)
	if (err != nil) {
		return nil, err
	}

	return &GAccount{
		Name: accountName,
		Service: service,
		Calendars: GCalendars{ All: calendars },
	}, nil
}

func (s *GAccount) SyncCalendars() (error) {
	calendars, err := getCalendars(s.Service)
	if (err != nil) {
		return fmt.Errorf("failed to sync calendars: %w", err)
	}

	s.Calendars.All = calendars
	return nil
}

func (s *GAccount) Init(serviceName string) (error) {
	service, err := getService(serviceName, s.Name)
	if (err != nil) {
		return err
	}

	s.Service = service
	return nil
}

func (s *GAccount) Events(calendarId string, filter *eventsfilter.EventsFilter) ([]*calendar.Event, error) {
	events, err := filter.Apply(s.Service.Events.List(calendarId)).Do()

	if err != nil {
		return nil, err
	}

	return events.Items, nil
}

func (s *GAccount) ResolveCalendars() ([]string) {
	if (len(s.Calendars.WhiteList) > 0) {
		return s.Calendars.WhiteList
	}

	return set.NewSet(s.Calendars.All...).Difference(set.NewSet(s.Calendars.BlackList...)).ToSlice();
}

func getService(serviceName string, accountName string) (*calendar.Service, error) {
	kr := typedkeyring.New[gaseed.GASeed](serviceName)
	gaSeed, err := kr.Load(accountName)
	if err != nil {
		return nil, err
	}

	client := gaSeed.GetClient()
	return calendar.NewService(context.Background(), option.WithHTTPClient(client))
}

func getCalendars(service *calendar.Service) ([]string, error) {
	calendars, err := service.CalendarList.List().Do()
	if err != nil {
		return nil, err
	}

	calendarNames := make([]string, len(calendars.Items))
	for i, cal := range calendars.Items {
		calendarNames[i] = cal.Id
	}

	return calendarNames, nil
}

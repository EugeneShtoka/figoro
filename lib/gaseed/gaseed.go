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
package gaseed

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)


type GASeed struct {
	Token 			*oauth2.Token
	Config 			*oauth2.Config
}

func New(clientID string, clientSecret string, bindAddress string, authEndpoint string) *GASeed {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL: fmt.Sprintf("http://%s%s", bindAddress, authEndpoint),
		Endpoint:     google.Endpoint,
		Scopes:       []string{calendar.CalendarReadonlyScope},
	}
	return &GASeed{	Config: config }
}

func (s *GASeed) SetToken(code string) (*GASeed, error) {
	var err error
	s.Token, err = s.Config.Exchange(context.Background(), code)
	return s, err
}

func (s *GASeed) GetClient() *http.Client {
	return s.Config.Client(context.Background(), s.Token)
}
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
package gauth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"text/template"

	"github.com/EugeneShtoka/figoro/lib/gaseed"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

var (
	host = "127.0.0.1"
	authEndpoint = "/codeResponse"
)

type GAServer struct {
	BindAddress		string
	AuthEndpoint	string
	State       	string
	Code			chan string
	Listener    	net.Listener
	Server      	*http.Server
	GASeed			*gaseed.GASeed
	Logger			*zerolog.Logger
}

type GAError struct {
	Name		string
	Description string
}

// New makes the webserver for collecting auth
func New(clientID string, clientSecret string, port string, logger *zerolog.Logger) *GAServer {
	bindAddress := fmt.Sprintf("%s:%s", host, port)
	return &GAServer{
		BindAddress: bindAddress,
		AuthEndpoint: authEndpoint,
		State:      uuid.New().String(),
		Logger:		logger,
		GASeed:		gaseed.New(clientID, clientSecret, bindAddress, authEndpoint),
		Code:		make(chan string, 1),
	}
}

	// Reply with the response to the user and to the channel
func (s *GAServer) reply(w http.ResponseWriter, res *GAError) {
	var (
		status int
		responseTemplate string
	)
	if (res == nil) {
		status = http.StatusOK
	} else {
		status = http.StatusBadRequest
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/html")
	var t = template.Must(template.New("authResponse").Parse(responseTemplate))
	var err = t.Execute(w, res)
	if  err != nil {
		s.Logger.Error().Msg(fmt.Sprintf("Could not execute template for web response: %s", err))
	}
}


func (s *GAServer) HandleAuth(w http.ResponseWriter, req *http.Request) {
	// Parse the form parameters and save them
	err := req.ParseForm()
	if err != nil {
		s.reply(w, &GAError{
			Name:        "Parse form error",
			Description: err.Error(),
		})
		return
	}

	// get code, error if empty
	var code = req.Form.Get("code")
	if code == "" {
		s.reply(w, &GAError{
			Name:        "Auth Error",
			Description: "No code returned by remote server",
		})
		return
	}

	// check state
	var state = req.Form.Get("state")
	if state != s.State {
		s.reply(w, &GAError{
			Name:        "Auth state doesn't match",
			Description: fmt.Sprintf("Expecting %q got %q", s.State, state),
		})
		return
	}

	// code OK
	s.reply(w, nil)
	s.Code <- req.FormValue("code")
}

// Init gets the internal web server ready to receive config details
func (gaServer *GAServer) Init() error {
	gaServer.Logger.Debug().Str("BindAddress", gaServer.BindAddress).Msg("Starting auth server")
	var mux = http.NewServeMux()
	gaServer.Server = &http.Server{
		Addr:    gaServer.BindAddress,
		Handler: mux,
	}
	gaServer.Server.SetKeepAlivesEnabled(false)

	mux.HandleFunc(gaServer.AuthEndpoint, gaServer.HandleAuth)

	var err error
	gaServer.Listener, err = net.Listen("tcp", gaServer.BindAddress)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	return nil
}

// Serve the auth server, doesn't return
func (s *GAServer) Serve() error {
	var err = s.Server.Serve(s.Listener)
	return fmt.Errorf("closed auth server with error: %w", err)
}

// Stop the auth server by closing its socket
func (s *GAServer) Stop() {
	s.Logger.Debug().Msg("Closing auth server")
	s.Listener.Close()
	s.Server.Close()
}

func (gaServer *GAServer) Authorize(ctx context.Context) (*gaseed.GASeed, error) {
	err := gaServer.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to start auth webserver: %w", err)
	}

	go gaServer.Serve()
	defer gaServer.Stop()

	// Open the URL for the user to visit
	var authUrl = gaServer.GASeed.Config.AuthCodeURL(gaServer.State, oauth2.AccessTypeOffline)
	open.Start(authUrl)

	fmt.Printf("Waiting for code\n")
	var code = <- gaServer.Code

	if	code == "" {
		return nil, fmt.Errorf("failed to start auth webserver: %w", err)
	} 

	return gaServer.GASeed.SetToken(code)
}


// Copyright 2020 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"github.com/urfave/cli"
)

const (
	programName = "bart-server"
)

// BartApp contains everything needed to run the BART demo client or server.
type BartApp struct {
	*cli.App
	address        string
	grpcAddress    string
	tlsCert        string
	tlsKey         string
	tlsDisable     bool
	modelPath      string
	output         string
	requestText    string
	requestText2   string
	commaSepLabels string
	multiClass     bool
}

// NewBartApp returns BartApp objects. The app can be used as both a client and a server.
func NewBartApp() *BartApp {
	app := &BartApp{
		App: cli.NewApp(),
	}
	app.Name = programName
	app.Usage = "A demo for sequence-classification based on BART."
	app.Commands = []cli.Command{
		newServerCommandFor(app),
		newClientCommandFor(app),
	}
	return app
}

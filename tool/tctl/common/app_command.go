/*
Copyright 2020 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"context"
	"os"

	"github.com/gravitational/kingpin"
	"github.com/gravitational/teleport/lib/auth"
	"github.com/gravitational/teleport/lib/defaults"
	"github.com/gravitational/teleport/lib/service"
	"github.com/gravitational/teleport/lib/services"
	"github.com/gravitational/trace"
)

// AppsCommand implements "tctl apps" group of commands.
type AppsCommand struct {
	config *service.Config

	// format is the output format, e.g. text or json.
	format string

	// appsList implements the "tctl apps ls" subcommand.
	appsList *kingpin.CmdClause
}

// Initialize allows AppsCommand to plug itself into the CLI parser
func (c *AppsCommand) Initialize(app *kingpin.Application, config *service.Config) {
	c.config = config

	apps := app.Command("apps", "Operate on applications registered with the cluster.")
	c.appsList = apps.Command("ls", "List all applications registered with the cluster.")
	c.appsList.Alias(ListAppsHelp)
}

// TryRun attempts to run subcommands like "apps ls".
func (c *AppsCommand) TryRun(cmd string, client auth.ClientI) (match bool, err error) {
	switch cmd {
	case c.appsList.FullCommand():
		err = c.ListApps(client)
	default:
		return false, nil
	}
	return true, trace.Wrap(err)
}

// ListApps prints the list of applications that have recently sent heartbeats
// to the cluster.
func (c *AppsCommand) ListApps(client auth.ClientI) error {
	apps, err := client.GetApps(context.TODO(), defaults.Namespace, services.SkipValidation())
	if err != nil {
		return trace.Wrap(err)
	}
	coll := &appCollection{apps: apps}
	coll.writeText(os.Stdout)
	return nil
}

package main

import (
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/plugin"
)

var (
	HOOK_ID     string
	ADMIN_TOKEN string
	TEAM_NAME   string
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()
	ADMIN_TOKEN = config.ADMIN_TOKEN
	TEAM_NAME = config.Clone().TEAM_NAME

	if r.Header.Get("Data") == "" {
		w.Write([]byte("Data Null"))
		return
	}

	switch r.URL.Path {
	default:
		http.NotFound(w, r)
	}
}

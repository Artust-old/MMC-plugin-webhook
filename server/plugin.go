package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/mattermost/mattermost-server/plugin"
)

var (
	HOOK_URL    string
	ADMIN_TOKEN string
	TEAM_NAME   string

	PORT_LISTEN string
	URL_SITE    string
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

type BodyGitHub struct {
	Action  string `json:"action"`
	Ref     string `json:"ref"`
	Zen     string `json:"zen"`
	Hook_id int    `json:"hook_id"`
	Text    string `json:"text"`
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()
	HOOK_URL = config.HOOK_URL
	ADMIN_TOKEN = config.ADMIN_TOKEN
	TEAM_NAME = config.Clone().TEAM_NAME

	bodyRequest, err := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyRequest))
	var dataGitHub BodyGitHub
	json.Unmarshal([]byte(bodyRequest), &dataGitHub)

	if err != nil {
		w.Write([]byte("Can't read body"))
		return
	}

	switch r.URL.Path {
	case "/github":
		if dataGitHub.Action != "" {
			p.GitHub(c, w, r, dataGitHub)
		}
	case "/status":
		p.handleStatus(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) handleStatus(w http.ResponseWriter, r *http.Request) {
	configuration := p.getConfiguration()

	var response = struct {
		Enabled bool `json:"enabled"`
	}{
		Enabled: !configuration.disabled,
	}

	responseJSON, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

func (p *Plugin) GitHub(c *plugin.Context, w http.ResponseWriter, r *http.Request, dataGitHub BodyGitHub) {
	pushForm := fmt.Sprintf(`{
		"text":"
		*Action:* %s
		*Branch name:* %s"
		}`,
		dataGitHub.Action, dataGitHub.Ref)

	req, err := http.NewRequest("POST", HOOK_URL, bytes.NewBufferString(pushForm))
	req.AddCookie(&http.Cookie{Name: "MMAUTHTOKEN", Value: ADMIN_TOKEN})
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New("\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	something, _ := ioutil.ReadAll(resp.Body)
	w.Write([]byte(something))
	return
}

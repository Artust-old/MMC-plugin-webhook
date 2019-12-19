package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
)

func (p *Plugin) postPluginMessage(channelID string, msg string) *model.AppError {
	configuration := p.getConfiguration()

	if configuration.disabled {
		return nil
	}

	msg = fmt.Sprintf("_%s_", msg)
	return &model.AppError{}
}

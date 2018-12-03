package got

import (
	"context"

	"github.com/nlopes/slack"
)

type Slack interface {
	AuthTestContext(context.Context) (*slack.AuthTestResponse, error)
	NewRTM(...slack.RTMOption) *slack.RTM
}

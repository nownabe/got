package got

import (
	"context"

	"github.com/nlopes/slack"
)

type mockSlack struct {
	err    error
	user   string
	userID string
}

func (c *mockSlack) AuthTestContext(ctx context.Context) (*slack.AuthTestResponse, error) {
	if c.err != nil {
		return nil, c.err
	}
	return &slack.AuthTestResponse{
		User:   c.user,
		UserID: c.userID,
	}, nil
}

func (c *mockSlack) NewRTM(options ...slack.RTMOption) *slack.RTM {
	return nil
}

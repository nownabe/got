package got

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test(t *testing.T) {
	ctx := context.Background()
	logger := zap.NewNop().Sugar()

	user := "user"
	userID := "userID"

	s := &mockSlack{
		err:    nil,
		user:   user,
		userID: userID,
	}

	g, err := NewWithSlack(ctx, s, WithLogger(logger))
	assert.Equal(t, user, g.Name())
	assert.Equal(t, userID, g.ID())
	assert.Nil(t, err)

	s = &mockSlack{err: errors.New("")}
	_, err = NewWithSlack(ctx, s, WithLogger(logger))
	assert.Error(t, err)
}

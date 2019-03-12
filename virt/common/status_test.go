package common

import (
	"testing"

	"github.com/projecteru2/yavirt/test/assert"
)

var allStatus = []string{
	StatusPending,
	StatusCreating,
	StatusStarting,
	StatusRunning,
	StatusStopping,
	StatusStopped,
	StatusMigrating,
	StatusResizing,
	StatusDestroying,
	StatusDestroyed,
}

func TestStatusForward(t *testing.T) {
	var cases = []struct {
		forward string
		allowed map[string]struct{}
	}{
		{
			StatusPending,
			allow([]string{StatusPending, ""}),
		},
		{
			StatusCreating,
			allow([]string{StatusCreating, StatusPending}),
		},
		{
			StatusStarting,
			allow([]string{StatusStarting, StatusCreating, StatusStopped, StatusRunning}),
		},
		{
			StatusRunning,
			allow([]string{StatusRunning, StatusStarting}),
		},
		{
			StatusStopping,
			allow([]string{StatusStopping, StatusRunning, StatusStopped}),
		},
		{
			StatusStopped,
			allow([]string{StatusStopped, StatusStopping, StatusMigrating}),
		},
		{
			StatusMigrating,
			allow([]string{StatusMigrating, StatusStopped, StatusResizing}),
		},
		{
			StatusResizing,
			allow([]string{StatusResizing, StatusStopped}),
		},
		{
			StatusDestroying,
			allow([]string{StatusDestroying, StatusStopped, StatusDestroyed}),
		},
		{
			StatusDestroyed,
			allow([]string{StatusDestroyed, StatusDestroying}),
		},
	}

	for _, c := range cases {
		var next = c.forward
		for _, now := range allStatus {
			if _, exists := c.allowed[now]; exists {
				assert.True(t, checkForward(now, next), "expect true of %s to %s", now, next)
			} else {
				assert.False(t, checkForward(now, next), "expect false of %s to %s", now, next)
			}
		}
	}
}

func allow(st []string) map[string]struct{} {
	var m = map[string]struct{}{}

	for _, elem := range st {
		m[elem] = struct{}{}
	}

	return m
}

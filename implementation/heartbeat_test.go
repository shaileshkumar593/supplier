package implementation_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeartbeat(t *testing.T) {
	var ctx = context.Background()
	resp, err := svc.HeartBeat(ctx)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "Connection alive", resp.Body)
}

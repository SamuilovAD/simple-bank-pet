package gapi

import (
	db "github.com/SamuilovAD/simple-bank-pet/db/sqlc"
	"github.com/SamuilovAD/simple-bank-pet/util"
	"github.com/SamuilovAD/simple-bank-pet/worker"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func newTestServer(t *testing.T, store db.Store, taskDistributer worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store, taskDistributer)
	require.NoError(t, err)

	return server
}

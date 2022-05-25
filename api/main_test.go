package api

import (
	"os"
	"testing"
	"time"

	db "github.com/Franklynoble/bankapp/db/sqlc"
	"github.com/Franklynoble/bankapp/db/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {

	config := util.Config{
		TokenSymmetricKey:  util.Randomstring(32),
		AccessTokenDuraion: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

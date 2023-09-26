package adapter

import (
	"testing"

	"github.com/go-playground/assert/v2"
	log "github.com/nxs23/utility-library-go/logging/logger"
)

func TestNewAdapter(t *testing.T) {
	conf := Config{}
	l := log.NewLogger("mysql-adapter")
	a := NewAdapter(conf, l)
	assert.Equal(t, a.IsConnected, false)
	assert.Equal(t, a.Config, conf)
	assert.Equal(t, a.Logger, l)
}

func TestDefaultAdapter(t *testing.T) {
	a := DefaultAdapter()

	assert.Equal(t, "mysql-adapter", a.Logger.GetLogger().Module)
}

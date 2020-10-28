package docker

import (
	"testing"
)

func TestStartContainer(t *testing.T) {
	cnd := StartContainer(t, "postgres:11.1-alpine")
	DumpContainerLogs(t, cnd)
	StopContainer(t, cnd)
}

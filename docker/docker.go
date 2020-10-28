package docker

import (
	"bytes"
	"encoding/json"
	"net"
	"os/exec"
	"testing"
)

type Container struct {
	ID   string
	Host string // IP:Port
}

func StartContainer(t *testing.T, image string) *Container {
	t.Helper()

	cmd := exec.Command("docker", "run", "-P", "-d", image)
	out := bytes.Buffer{}
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("executing creating docker container %#v", err)
	}

	id := out.String()[0:12]
	t.Log("DB container id", id)

	out.Reset()

	cmd = exec.Command("docker", "inspect", id)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		t.Log("couldn't inspect docker container", err)
	}

	var doc []struct {
		NetworkSettings struct {
			Ports struct {
				TCP5432 []struct {
					HostIp   string `json:"HostIp"`
					HostPort string `json:"HostPort"`
				} `json:"5432/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}

	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("couldn't decode json %#v", err)
	}

	host := doc[0].NetworkSettings.Ports.TCP5432[0].HostIp
	port := doc[0].NetworkSettings.Ports.TCP5432[0].HostPort

	c := &Container{
		ID:   id,
		Host: net.JoinHostPort(host, port),
	}

	t.Log("DB Host:", c.Host)
	return c
}

func StopContainer(t *testing.T, c *Container) {
	t.Helper()

	if err := exec.Command("docker", "stop", c.ID).Run(); err != nil {
		t.Fatalf("couldn't stop  container %v", err)
	}

	if err := exec.Command("docker", "rm", c.ID).Run(); err != nil {
		t.Fatalf("couldn't remove container %v", err)
	}

	t.Log("removed container", c.ID)

}

func DumpContainerLogs(t *testing.T, c *Container) {
	t.Helper()
	out, err := exec.Command("docker", "logs", c.ID).CombinedOutput()
	if err != nil {
		t.Fatalf("fetching container logs %#v", err)
	}
	t.Logf("Logs %s\n%s", c.ID, out)
}

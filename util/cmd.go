package util

import (
	"io"
	"os"
	"os/exec"
)

type Cmd struct {
	executable string
	mac        string
	sudo       string
	output     io.Writer
}

func NewCmd(out io.Writer, MAC string) *Cmd {
	if out == nil {
		out = os.Stdout
	}
	return &Cmd{
		executable: "bluetoothctl",
		mac:        MAC,
		sudo:       "sudo",
		output:     out,
	}
}

func (c *Cmd) run(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = c.output
	cmd.Stderr = c.output
	return cmd.Run()
}

func (c *Cmd) Connect() error {
	return c.run(c.executable, "connect", c.mac)
}

func (c *Cmd) Disconnect() error {
	return c.run(c.executable, "disconnect")
}

func (c *Cmd) Poweron() error {
	return c.run(c.executable, "power", "on")
}

func (c *Cmd) Poweroff() error {
	return c.run(c.executable, "power", "off")
}

func (c *Cmd) RestartService() error {
	return c.run(c.sudo, "systemctl", "restart", "bluetooth")
}

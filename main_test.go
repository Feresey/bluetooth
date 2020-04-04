package main

import (
	"context"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testingMAC     = "00:11:22:33:44:55:66"
	sleep          = 10 * time.Millisecond
	scanCommand    = "bluetoothctl scan on"
	deviceFound    = "[NEW] Device " + testingMAC
	deviceNotFound = "[NOT NEW] Device"
)

func debugCmd(c *cmd, mac string) *cmd {
	c.execCommand = helperCommand
	c.execCommandContext = helperCommandContext
	c.MAC = mac
	c.Quiet = true
	return c
}

func TestCmdScanByInterval(t *testing.T) {
	c := debugCmd(newCmd(), testingMAC)

	t.Run("on first hit", func(t *testing.T) {
		result[scanCommand] = &trapExec{
			text: deviceFound,
		}

		ok := c.scanByInterval(sleep)
		assert.Equal(t, ok, true)
	})

	t.Run("on ten hit", func(t *testing.T) {
		fake := &trapExec{
			text: deviceNotFound,
		}
		result[scanCommand] = fake

		for i := 0; i < 10; i++ {
			ok := c.scanByInterval(sleep)
			assert.Equal(t, ok, false)
		}

		assert.Equal(t, 10, fake.called)

		fake.text = deviceFound

		ok := c.scanByInterval(sleep)
		assert.Equal(t, ok, true)
	})
}

func TestCmdScan(t *testing.T) {
	c := debugCmd(newCmd(), testingMAC)

	t.Run("on first hit", func(t *testing.T) {
		result[scanCommand] = &trapExec{
			text: deviceFound,
		}

		var (
			ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(time.Second))
			done        = make(chan struct{})
		)
		defer cancel()

		go func() {
			c.scan(context.TODO(), sleep)
			done <- struct{}{}
		}()

		select {
		case <-ctx.Done():
			t.Error("Too long")
		case <-done:
		}
	})

	t.Run("on tenth hit", func(t *testing.T) {
		fake := &trapExec{
			text: deviceNotFound,
		}
		result[scanCommand] = fake

		var (
			ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(time.Second))
			done        = make(chan struct{})
		)
		defer cancel()

		go func() {
			c.scan(context.TODO(), sleep)
			done <- struct{}{}
		}()

		go func() {
			for fake.called < 10 {
				time.Sleep(sleep)
			}

			fake.text = deviceFound
		}()

		select {
		case <-ctx.Done():
			t.Error("Too long")
		case <-done:
		}
	})
}

func TestCmd_pair(t *testing.T) {
	type fields struct {
		execCommand func(name string, args ...string) *exec.Cmd
		executable  string
		mac         string
		sudo        string
	}
	type args struct {
		ctx        context.Context
		cancelScan func()
		sleep      time.Duration
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantInfo string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cmd{
				execCommand: tt.fields.execCommand,
				executable:  tt.fields.executable,
				MAC:         tt.fields.mac,
				sudo:        tt.fields.sudo,
			}
			gotInfo, err := c.Pair(tt.args.ctx, tt.args.cancelScan, tt.args.sleep)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.pair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInfo != tt.wantInfo {
				t.Errorf("Cmd.pair() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func TestCmd_connect(t *testing.T) {
	type fields struct {
		execCommand func(name string, args ...string) *exec.Cmd
		executable  string
		mac         string
		sudo        string
	}
	tests := []struct {
		name     string
		fields   fields
		wantInfo string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cmd{
				execCommand: tt.fields.execCommand,
				executable:  tt.fields.executable,
				MAC:         tt.fields.mac,
				sudo:        tt.fields.sudo,
			}
			gotInfo, err := c.Connect()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInfo != tt.wantInfo {
				t.Errorf("Cmd.connect() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func TestCmd_disconnect(t *testing.T) {
	type fields struct {
		execCommand func(name string, args ...string) *exec.Cmd
		executable  string
		mac         string
		sudo        string
	}
	tests := []struct {
		name     string
		fields   fields
		wantInfo string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cmd{
				execCommand: tt.fields.execCommand,
				executable:  tt.fields.executable,
				MAC:         tt.fields.mac,
				sudo:        tt.fields.sudo,
			}
			gotInfo, err := c.Disconnect()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.disconnect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInfo != tt.wantInfo {
				t.Errorf("Cmd.disconnect() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func TestCmd_on(t *testing.T) {
	type fields struct {
		execCommand func(name string, args ...string) *exec.Cmd
		executable  string
		mac         string
		sudo        string
	}
	tests := []struct {
		name     string
		fields   fields
		wantInfo string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cmd{
				execCommand: tt.fields.execCommand,
				executable:  tt.fields.executable,
				MAC:         tt.fields.mac,
				sudo:        tt.fields.sudo,
			}
			gotInfo, err := c.On()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.on() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInfo != tt.wantInfo {
				t.Errorf("Cmd.on() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func TestCmd_off(t *testing.T) {
	type fields struct {
		execCommand func(name string, args ...string) *exec.Cmd
		executable  string
		mac         string
		sudo        string
	}
	tests := []struct {
		name     string
		fields   fields
		wantInfo string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cmd{
				execCommand: tt.fields.execCommand,
				executable:  tt.fields.executable,
				MAC:         tt.fields.mac,
				sudo:        tt.fields.sudo,
			}
			gotInfo, err := c.Off()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.off() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInfo != tt.wantInfo {
				t.Errorf("Cmd.off() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func TestCmd_restart(t *testing.T) {
	type fields struct {
		execCommand func(name string, args ...string) *exec.Cmd
		executable  string
		mac         string
		sudo        string
	}
	tests := []struct {
		name     string
		fields   fields
		wantInfo string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cmd{
				execCommand: tt.fields.execCommand,
				executable:  tt.fields.executable,
				MAC:         tt.fields.mac,
				sudo:        tt.fields.sudo,
			}
			gotInfo, err := c.Restart()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.restart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInfo != tt.wantInfo {
				t.Errorf("Cmd.restart() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func TestCmd_remove(t *testing.T) {
	type fields struct {
		execCommand func(name string, args ...string) *exec.Cmd
		executable  string
		mac         string
		sudo        string
	}
	tests := []struct {
		name     string
		fields   fields
		wantInfo string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cmd{
				execCommand: tt.fields.execCommand,
				executable:  tt.fields.executable,
				MAC:         tt.fields.mac,
				sudo:        tt.fields.sudo,
			}
			gotInfo, err := c.Remove()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInfo != tt.wantInfo {
				t.Errorf("Cmd.remove() = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

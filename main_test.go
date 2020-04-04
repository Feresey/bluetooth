package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testingMAC     = "00:11:22:33:44:55:66"
	sleep          = 5 * time.Millisecond
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
			c.scan(ctx, sleep)
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
			c.scan(ctx, sleep)
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

	t.Run("cancel", func(t *testing.T) {
		result[scanCommand] = &trapExec{
			text: deviceNotFound,
		}

		var (
			localSleep  = 100 * time.Millisecond
			ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(localSleep))
			done        = make(chan struct{})
		)

		go func() {
			c.scan(ctx, sleep)
			done <- struct{}{}
		}()

		time.Sleep(localSleep / 2)
		cancel()

		select {
		case <-time.Tick(localSleep / 2):
			t.Error("Too long")
		case <-done:
		}
	})
}

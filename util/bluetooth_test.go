package  util

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func TestCmd(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	c := &Cmd{
		executable: "echo",
		sudo:       "echo",
		output:     buffer,
	}
	tests := []struct {
		name           string
		MAC            string
		expectedOutput []byte
		method         func() error
	}{
		{
			name:           "connect",
			MAC:            "MAC",
			expectedOutput: []byte("connect MAC\n"),
			method:         func() error { return c.Connect("MAC") },
		},
		{
			name:           "disconnect",
			expectedOutput: []byte("disconnect\n"),
			method:         c.Disconnect,
		},
		{
			name:           "poweron",
			expectedOutput: []byte("power on\n"),
			method:         c.Poweron,
		},
		{
			name:           "poweroff",
			expectedOutput: []byte("power off\n"),
			method:         c.Poweroff,
		},
		{
			name:           "restart",
			expectedOutput: []byte("systemctl restart bluetooth\n"),
			method:         c.RestartService,
		},
	}

	t.Run("new", func(t *testing.T) {
		if !reflect.DeepEqual(NewCmd(nil), &Cmd{executable: "bluetoothctl", sudo: "sudo", output: os.Stdout}) {
			t.Errorf("Struct does not match")
		}
	})

	for _, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {
			if err := tt.method(); err != nil {
				t.Errorf("cmd.connect() error = %v", err)
			}
			if !bytes.Equal(tt.expectedOutput, buffer.Bytes()) {
				t.Errorf("Wrong output.\nGiven:\n%q\nExpected:\n%q", buffer, tt.expectedOutput)
			}
			buffer.Reset()
		})
	}
}

package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var mu sync.Mutex

type cmd struct {
	execCommand        func(name string, args ...string) *exec.Cmd
	execCommandContext func(ctx context.Context, name string, args ...string) *exec.Cmd

	executable string
	sudo       string
	MAC        string
	Quiet      bool
}

func newCmd() *cmd {
	return &cmd{
		execCommand:        exec.Command,
		execCommandContext: exec.CommandContext,

		executable: "bluetoothctl",
		sudo:       "xfsudo",
	}
}

func main() {
	command := newCmd()

	flag.Usage = func() {
		fmt.Println(`This program parse first argument by characters:
		"r": restart bluetooth daemon
		"+": connect to the MAC
		"-": disconnect from the MAC
		"d": remove device and connect cleanly`)
	}

	flag.StringVar(&command.MAC, "mac", "5C:FB:7C:77:23:E2", "Address of bluetooth device")
	flag.BoolVar(&command.Quiet, "q", false, "Do not verbose output")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		gracefullShutdown()
		cancel()
	}()

	for _, elem := range flag.Args()[0] {
		switch elem {
		case '+':
			must(command.On())
			must(command.Connect())
		case '-':
			must(command.Disconnect())
			must(command.Off())
		case 'r':
			must(command.Restart())
		case 'c':
			must(command.On())
			must(command.Remove())

			must(command.Scan(ctx, 3*time.Second))
			must(command.Pair(ctx, cancel, time.Second))
			must(command.Connect())
		default:
			log.Printf("No such command: %v", elem)
		}
	}
}

func gracefullShutdown() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
	fmt.Println("Caught interrupt signal, stopping all processes")
}

func must(info string, err error) (ok bool) {
	msg := new(strings.Builder)

	msg.WriteString(info)
	msg.WriteString("...\t: ")

	mu.Lock()
	if err != nil {
		log.SetLevel(log.ErrorLevel)
		msg.WriteString(err.Error())
		ok = false
	} else {
		log.SetLevel(log.InfoLevel)
		msg.WriteString("SUCCESS")
		ok = true
	}
	mu.Unlock()

	log.Print(msg.String())
	return
}

func (c *cmd) execHere(command string, args ...string) *exec.Cmd {
	res := c.execCommand(command, args...)

	if !c.Quiet {
		res.Stdout = os.Stdout
		res.Stderr = os.Stderr
	}

	return res
}

func (c *cmd) Scan(ctx context.Context, active time.Duration) (info string, err error) {
	go c.scan(ctx, active)

	return "Scan avaliable devices", nil
}

func (c *cmd) scan(ctx context.Context, active time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if c.scanByInterval(active) {
				log.Info("Device found")
				return
			}
		}
	}
}

func (c *cmd) scanByInterval(active time.Duration) bool {
	var (
		ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(active))
		cmd         = c.execCommandContext(ctx, c.executable, "scan", "on")
		waitFor     = fmt.Sprintf(`[NEW] Device %s`, c.MAC)
	)
	defer cancel()

	// ошибка будет всегда, т.к. я убиваю нахер процесс
	out, _ := cmd.Output()

	if !c.Quiet {
		_, err := os.Stdout.Write(out)
		if err != nil {
			log.Warn(err)
		}
	}

	s := bufio.NewScanner(bytes.NewReader(out))

	for s.Scan() {
		if strings.HasPrefix(s.Text(), waitFor) {
			return true
		}
	}

	return false
}

func (c *cmd) Pair(ctx context.Context, cancelScan func(), sleep time.Duration) (info string, err error) {
	info = "Pair with specified device"

loop:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err = c.execHere(c.executable, "pair", c.MAC).Run()
			must(info, err)
			if err == nil {
				break loop
			} else {
				time.Sleep(sleep)
			}
		}
	}

	cancelScan()

	return
}

func (c *cmd) Connect() (info string, err error) {
	return "Connect to specified device",
		c.execHere(c.executable, "connect", c.MAC).Run()
}

func (c *cmd) Disconnect() (info string, err error) {
	return "Disconnect specified device",
		c.execHere(c.executable, "disconnect").Run()
}

func (c *cmd) On() (info string, err error) {
	return "Power on bluetooth adapter",
		c.execHere(c.executable, "power", "on").Run()
}

func (c *cmd) Off() (info string, err error) {
	return "Power off bluetooth adapter",
		c.execHere(c.executable, "power", "off").Run()
}

func (c *cmd) Restart() (info string, err error) {
	return "Restart bluetooth service",
		c.execHere(c.sudo, "systemctl", "restart", "bluetooth").Run()
}

func (c *cmd) Remove() (info string, err error) {
	return "Remove specified device",
		c.execHere(c.executable, "remove", c.MAC).Run()
}

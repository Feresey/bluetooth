package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

type Cmd struct {
	executable string
	mac        string
	sudo       string
}

func newCmd(mac string) *Cmd {
	return &Cmd{
		executable: "bluetoothctl",
		mac:        mac,
		sudo:       "xfsudo",
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println(`This program parse first argument by characters:
		"r": restart bluetooth daemon and reconnect speifier device
		"c": connect to the MAC
		"d": disconnect
		"+": power on bluetooth
		"-": power off bluetooth`)
	}

	var mac string

	flag.StringVar(&mac, "mac", "5C:FB:7C:77:23:E2", "MAC address of bluetooth device")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		return
	}

	command := newCmd(mac)

	go gracefullShutdown()

	for _, elem := range flag.Args()[0] {
		switch elem {
		case '+':
			must(command.on())
			must(command.connect())
		case '-':
			must(command.disconnect())
			must(command.off())
		case 'r':
			must(command.restart())
		case 'd':
			must(command.on())
			must(command.remove())
			must(command.scan())
			must(command.pair())
			must(command.connect())
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

func must(info string, err error) {
	msg := new(strings.Builder)

	msg.WriteString(info)
	msg.WriteString("...\t: ")

	if err != nil {
		msg.WriteString(err.Error())
	} else {
		msg.WriteString("SUCCESS")
	}

	fmt.Println(msg.String())
}

func execHere(command string, args ...string) *exec.Cmd {
	res := exec.Command(command, args...)

	res.Stdout = os.Stdout
	res.Stderr = os.Stderr

	return res
}

func (c *Cmd) connect() (info string, err error) {
	return "Connect to specified device",
		execHere(c.executable, "connect", c.mac).Run()
}

func (c *Cmd) disconnect() (info string, err error) {
	return "Disconnect specified device",
		execHere(c.executable, "disconnect").Run()
}

func (c *Cmd) on() (info string, err error) {
	return "Power on bluetooth adapter",
		execHere(c.executable, "power", "on").Run()
}

func (c *Cmd) off() (info string, err error) {
	return "Power off bluetooth adapter",
		execHere(c.executable, "power", "off").Run()
}

func (c *Cmd) restart() (info string, err error) {
	return "Restart bluetooth service",
		execHere(c.sudo, "systemctl", "restart", "bluetooth").Run()
}

func (c *Cmd) remove() (info string, err error) {
	return "Remove specified device",
		execHere(c.executable, "remove", c.mac).Run()
}

func (c *Cmd) scan() (info string, err error) {
	info = "Scan avaliable devices"

	var (
		ctx, cancel = context.WithCancel(context.Background())
		cmd         = exec.CommandContext(ctx, c.executable, "scan", "on")
		waitFor     = fmt.Sprintf(`[NEW] Device %s`, c.mac)
	)

	stop := func() {
		cancel()
		must("Cancelling scan", cmd.Wait())
	}

	go func() {
		time.Sleep(5 * time.Second)
		stop()
	}()
	err = cmd.Start()

	out, er := cmd.Output()
	if er != nil {
		err = er
		return
	}

	s := bufio.NewScanner(bytes.NewReader(out))

	for s.Scan() {
		if strings.HasPrefix(s.Text(), waitFor) {
			return
		}
	}

	return
}

func (c *Cmd) pair() (info string, err error) {
	return "Pair with specified device",
		execHere(c.executable, "pair", c.mac).Run()
}

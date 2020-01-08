package main

import (
	"flag"
	"fmt"
	"github.com/Feresey/bluetooth/util"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Println(`This program parse first argument by characters:
		"r": restart bluetooth daemon
		"c": connect to the MAC
		"d": disconnect
		"+": power on bluetooth
		"-": power off bluetooth`)
	}
	MAC := flag.String("a", "5C:FB:7C:77:23:E2", "MAC adress of bluetooth device")
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		return
	}
	command := util.NewCmd(nil, *MAC)

	for _, elem := range os.Args[1] {
		switch elem {
		case 'c':
			if err := command.Poweron(); err != nil {
				panic(err)
			}
			if err := command.Connect(); err != nil {
				panic(err)
			}
		case 'd':
			if err := command.Disconnect(); err != nil {
				panic(err)
			}
		case '-':
			if err := command.Poweroff(); err != nil {
				panic(err)
			}
		case '+':
			if err := command.Poweron(); err != nil {
				panic(err)
			}
		case 'r':
			if err := command.RestartService(); err != nil {
				panic(err)
			}
		}

	}
}

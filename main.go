package main

import (
	"flag"
	"fmt"
	"log"
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

	MAC := flag.String("mac", "5C:FB:7C:77:23:E2", "MAC adress of bluetooth device")

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	command := newCmd(nil, *MAC)

	for _, elem := range os.Args[1] {
		switch elem {
		case 'c':
			if err := command.Poweron(); err != nil {
				log.Fatal(err)
			}
			if err := command.Connect(); err != nil {
				log.Fatal(err)
			}
		case 'd':
			if err := command.Disconnect(); err != nil {
				log.Fatal(err)
			}
		case '-':
			if err := command.Poweroff(); err != nil {
				log.Fatal(err)
			}
		case '+':
			if err := command.Poweron(); err != nil {
				log.Fatal(err)
			}
		case 'r':
			if err := command.RestartService(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

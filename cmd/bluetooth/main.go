package main

import (
	"github.com/Feresey/bluetooth/util"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Need one more argument")
	}
	command := util.NewCmd(nil)
	MAC := "5C:FB:7C:77:23:E2"

	for _, elem := range os.Args[1] {
		switch elem {
		case 'c':
			if err := command.Poweron(); err != nil {
				panic(err)
			}
			if err := command.Connect(MAC); err != nil {
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

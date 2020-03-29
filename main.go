package main

import (
	"fmt"

	"github.com/latipovsharif/bot-manager/platforms"
)

var runnableArr []platforms.Platform

func init() {
	runnableArr = []platforms.Platform{
		//&platforms.Telegram{},
		&platforms.Whatsapp{},
	}
}

func main() {
	for _, r := range runnableArr {
		if err := r.Run(); err != nil {
			panic(fmt.Sprintf("cannot run platform %T, err: %v", r, err))
		}
	}
}

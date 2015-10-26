package main

import (
	"io/ioutil"
	"log"

	"github.com/asteris-llc/consul-dynamic/commands"
)

const Name = "consul-dynamic"
const Version = "0.0.0"

func main() {
	log.SetOutput(ioutil.Discard)
	commands.Run()
}

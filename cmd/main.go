package main

import (
	"log"

	"mygeth"
)

func main() {
	node := mygeth.MakeFullNode()
	if err := node.Start(); err != nil {
		log.Fatal("node start failed %s", err.Error())
	}
	node.Wait()
}

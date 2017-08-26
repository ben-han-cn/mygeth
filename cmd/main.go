package main

import (
	"flag"
	"log"

	"mygeth"
)

var (
	dataDir string
)

func init() {
	flag.StringVar(&dataDir, "d", "", "data dir")
}

func main() {
	node := mygeth.MakeFullNode(dataDir)
	if err := node.Start(); err != nil {
		log.Fatal("node start failed %s", err.Error())
	}
	node.Wait()
}

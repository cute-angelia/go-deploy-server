package main

import (
	"flag"
	"go-deploy/internal/helper"
	"log"
)

func main() {
	cmdstr := flag.String("cmd", "", "")
	flag.Parse()
	rest, err := helper.RunShellCmd(*cmdstr)

	log.Println("rest:", string(rest))
	log.Println("err:", err)
}

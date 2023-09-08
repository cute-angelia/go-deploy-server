package main

import (
	"flag"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

func logErr(err error) {
	logging.Error(err)
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		network string
		addr    string
	)

	// Example command: go run client.go --network tcp --address ":9000" --concurrency 100 --packet_size 1024 --packet_batch 20 --packet_count 1000
	flag.StringVar(&network, "network", "tcp", "--network tcp")
	flag.StringVar(&addr, "address", "127.0.0.1:9000", "--address 127.0.0.1:9000")
	flag.Parse()

	gnet.NewClient()

}

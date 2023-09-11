package main

import (
	"flag"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"go-deploy/pkg/protocol"
	"log"
	"time"
)

type simpleClient struct {
	gnet.BuiltinEventEngine
	eng       gnet.Engine
	network   string
	addr      string
	multicore bool
}

func (s *simpleClient) OnBoot(eng gnet.Engine) (action gnet.Action) {
	log.Printf("running server on %s with multi-core=%t",
		fmt.Sprintf("%s://%s", s.network, s.addr), s.multicore)
	s.eng = eng
	return
}

func (s *simpleClient) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Println("OnOpen")

	log.Println(c.LocalAddr())
	log.Println(c.RemoteAddr())

	c.SetContext(new(protocol.SimpleCodec))
	//out = []byte("sweetness client connect")
	return
}

func (s *simpleClient) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	log.Println("OnClose")
	if err != nil {
		log.Printf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	return
}

func (s *simpleClient) OnTraffic(c gnet.Conn) (action gnet.Action) {
	log.Println("OnTraffic")
	codec := c.Context().(*protocol.SimpleCodec)
	var packets [][]byte
	for {
		msgType, data, err := codec.Decode(c)
		if err == protocol.ErrIncompletePacket {
			break
		}
		if err != nil {
			log.Printf("invalid packet: %v", err)
			return gnet.Close
		}

		log.Println(msgType)

		packet, _ := codec.Encode(protocol.MsgTypeData, data)
		packets = append(packets, packet)
	}
	//if n := len(packets); n > 1 {
	//	_, _ = c.Writev(packets)
	//} else if n == 1 {
	//	_, _ = c.Write(packets[0])
	//}
	return
}

func main() {
	var (
		addr string
	)

	log.SetFlags(log.Ltime | log.Lshortfile)

	flag.StringVar(&addr, "address", "127.0.0.1:9000", "--address 127.0.0.1:9000")
	flag.Parse()

	port := fmt.Sprintf(":%d", 0)
	client := &simpleClient{
		network: "tcp",
		addr:    port,
	}
	iclient, _ := gnet.NewClient(client, gnet.WithMulticore(false))
	cc, err := iclient.Dial("tcp", addr)
	if err := iclient.Start(); err != nil {
		log.Println(err)
	}

	codec := protocol.SimpleCodec{}

	// 心跳
	if err == nil {
		go func() {
			heartByte, _ := codec.Encode(protocol.MsgTypeHeart, nil)
			d := time.NewTicker(5 * time.Second)
			for {
				select {
				case <-d.C:
					cc.Write(heartByte)
				}
			}
		}()
	}

	for i := 0; i < 1000; i++ {
		data, _ := codec.Encode(protocol.MsgTypeData, []byte("hello"))
		if _, err := cc.Write(data); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 5)
	}

	select {}
}

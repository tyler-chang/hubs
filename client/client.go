package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tyler-chang/hubs"
)

func main() {
	// 解析服务器
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	// 连接服务器
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	p := &protocol.Protocol{}

	log.Println(p)

	// ping <--> pong
	for i := 0; i < 3; i++ {
		// write
		conn.Write(protocol.NewPacket([]byte("hello"), false).Serialize())

		// read
		p, err := p.ReadPacket(conn)
		if err == nil {
			packet := p.(*protocol.Packet)
			fmt.Printf("Server reply:[%v] [%v]\n", packet.GetLength(), string(packet.GetBody()))
		}

		time.Sleep(2 * time.Second)
	}

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

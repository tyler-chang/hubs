package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gansidui/gotcp"
	"github.com/tyler-chang/hubs"
)

// Callback 回调
type Callback struct{}

// OnConnect 新用户连接回调
func (pc *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	log.Println("OnConnect:", addr)
	return true
}

// OnMessage 接收到消息回调
func (pc *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*protocol.Packet)
	fmt.Printf("OnMessage:[%v] [%v]\n", packet.GetLength(), string(packet.GetBody()))
	c.AsyncWritePacket(protocol.NewPacket(packet.Serialize(), true), time.Second)
	return true
}

// OnClose 连接关闭回调
func (pc *Callback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose:", c.GetExtraData())
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &Callback{}, &protocol.Protocol{})

	// starts service
	go srv.Start(listener, time.Second)
	fmt.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

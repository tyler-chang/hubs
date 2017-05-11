package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"sync"

	"github.com/tyler-chang/gotcp"
)

var wg sync.WaitGroup
var maxOnlineClients int
var onlineClients int
var counterChan chan int

// Callback 事件回调对象
type Callback struct{}

// OnConnect 新用户接入事件回调
func (cb *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)
	// fmt.Println("OnConnect:", addr)
	counterChan <- 1
	return true
}

// OnMessage 新消息事件回调
func (cb *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	hp := p.(*gotcp.Hj212Packet)
	// fmt.Printf("OnMessage:[%v] [%v]\n", hp.GetLength(), string(hp.GetBody()))
	// 检查客户端发送的消息
	if string(hp.GetBody()) != "hello" {
		fmt.Fatal("客户端发送的值错误")
	}
	c.AsyncWritePacket(gotcp.BuildPacket([]byte("world")), time.Second)
	return true
}

// OnClose 用户连接关闭事件回调
func (cb *Callback) OnClose(c *gotcp.Conn) {
	// fmt.Println("OnClose:", c.GetExtraData())
	counterChan <- -1
}

func counter() {
	wg.Add(1)
	defer wg.Done()

	for {
		n, ok := <-counterChan
		if ok == false {
			return
		}
		onlineClients += n
		if onlineClients > maxOnlineClients {
			maxOnlineClients = onlineClients
			fmt.Println("maxOnlineClients:", maxOnlineClients)
		}
		fmt.Println("onlineClients:", onlineClients)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fatal(err)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	counterChan = make(chan int)

	go counter()

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
	srv := gotcp.NewServer(config, &Callback{}, &gotcp.Hj212Protocol{})

	// starts service
	go srv.Start(listener, time.Second)
	fmt.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	close(counterChan)
	wg.Wait()

	// stops service
	srv.Stop()
}

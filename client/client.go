package main

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sync"

	"github.com/tyler-chang/gotcp"
)

var wg sync.WaitGroup
var exitChan chan struct{}

func client() {
	wg.Add(1)
	defer wg.Done()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:8989")
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer conn.Close()

	protocol := &gotcp.Hj212Protocol{}
	// ping <--> pong
	for {
		select {
		case <-exitChan:
			return
		default:
		}

		// write
		packet := gotcp.BuildPacket([]byte("hello"))
		// fmt.Printf("Client send [%v] [%v]\n", packet.GetLength(), string(packet.GetBody()))
		conn.Write(packet.Serialize())

		// read
		p, err := protocol.ReadPacket(conn)
		if err == nil {
			packet := p.(*gotcp.Hj212Packet)
			// 检查服务器的回复
			if string(packet.GetBody()) != "world" {
				log.Fatal(errors.New("服务器返回值错误"))
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	wg = sync.WaitGroup{}
	exitChan = make(chan struct{})

	for i := 0; i < 60000; i++ {
		go client()
		time.Sleep(time.Millisecond * 2)
	}

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal: ", <-chSig)

	close(exitChan)
	wg.Wait()
}

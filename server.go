package main

import (
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
	"github.com/tyler-chang/gotcp"
)

var client *redis.Client
var wg sync.WaitGroup
var maxOnlineClients int
var onlineClients int
var counterChan chan int

// Callback 事件回调对象
type Callback struct{}

// OnConnect 新用户接入事件回调
func (cb *Callback) OnConnect(c *gotcp.Conn) bool {
	// 获取客户端 IP
	addr := c.GetRawConn().RemoteAddr()
	// 将客户端 IP 存入附加数据
	c.PutExtraData(addr)
	// 记录日志
	log.WithFields(log.Fields{
		"address": addr,
	}).Info("Client is connected")
	// 向计数器通道发送 +1 信号
	counterChan <- 1
	return true
}

// OnMessage 新消息事件回调
func (cb *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	// 断言类型
	hp := p.(*gotcp.Hj212Packet)

	// 记录日志
	log.WithFields(log.Fields{
		"length":  hp.GetLength(),
		"message": string(hp.GetBody()),
	}).Info("Accept a new message")

	// 检查客户端发送的消息
	if string(hp.GetBody()) != "hello" {
		log.Fatal("客户端发送的值错误")
	}

	// 向客户端发送数据
	c.AsyncWritePacket(gotcp.BuildPacket([]byte("world")), time.Second)
	return true
}

// OnClose 用户连接关闭事件回调
func (cb *Callback) OnClose(c *gotcp.Conn) {
	// 记录日志
	log.WithFields(log.Fields{
		"address": c.GetExtraData(),
	}).Info("Client is disconnected")
	// 向计数器通道发送 -1 信号
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
			// 记录最大在线连接数
			client.Set("maxOnlineClients", maxOnlineClients, 0)
			// log.WithFields(log.Fields{
			// 	"maxOnlineCount": maxOnlineClients,
			// }).Info("Max online clients count is changed")
		}
		// 记录当前在线连接数
		client.Set("onlineClients", onlineClients, 0)
		// log.WithFields(log.Fields{
		// 	"onlineCount": onlineClients,
		// }).Info("Online clients count is changed")
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       5,  // use default DB
	})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	file, _ := os.OpenFile("./log.log", os.O_CREATE|os.O_RDWR, 0755)
	defer func() {
		file.Close()
	}()
	log.SetOutput(file)

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
	// 记录开始日志
	log.WithFields(log.Fields{
		"address": listener.Addr(),
	}).Info("Server started")

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	// 记录系统信号
	log.WithFields(log.Fields{
		"signal": <-chSig,
	}).Info("Received system signal")

	close(counterChan)
	wg.Wait()

	// stops service
	srv.Stop()
}

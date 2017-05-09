package main

import (
	"log"
	"net"
	"time"
)

func ping() {
	log.Println("start ping...")

	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Println("dial error:", err)
		return
	}

	defer conn.Close()

	data := "12345"

	for i := 0; i < 10; i++ {
		conn.Write([]byte(data))
	}

	log.Println("end ping...")

}

func main() {
	for i := 0; i < 20000; i++ {
		go ping()
	}

	time.Sleep(time.Second * 1000)
}

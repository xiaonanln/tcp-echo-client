package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
)

var args struct {
	serverAddr string
	numClients int
}

var (
	sendData []byte
	counter  int64
)

func init() {
	sendData = make([]byte, 1024)
	rand.Read(sendData)
	go func() {
		for {
			time.Sleep(time.Second)
			counter := atomic.SwapInt64(&counter, 0)
			println(counter)
		}
	}()
}

func main() {
	parseArgs()

	for i := 0; i < args.numClients; i++ {
		newEchoClient(args.serverAddr)
	}

	select {}
}

func parseArgs() {
	flag.StringVar(&args.serverAddr, "server", "localhost:1234", "server address")
	flag.IntVar(&args.numClients, "N", 1, "number of clients")
	flag.Parse()
}

type EchoClient struct {
	serverAddr string
}

func (c *EchoClient) routine() {
reconnect:
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		log.Printf("connect failed")
		time.Sleep(time.Second)
		goto reconnect
	}

	log.Printf("connected: %s", conn.LocalAddr())
	recvBuffer := make([]byte, len(sendData))
	for {
		conn.Write(sendData)
		_, err := io.ReadFull(conn, recvBuffer)
		if err != nil {
			break
		}
		atomic.AddInt64(&counter, 1)
	}

	goto reconnect
}

func newEchoClient(serverAddr string) *EchoClient {
	client := &EchoClient{
		serverAddr: serverAddr,
	}
	go client.routine()
	return client
}

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/KaiserWerk/Loggy/internal/logging"
)

const bufferSize = 1024

func main() {
	bl, err := logging.New(".") // TODO: use flag or env var
	if err != nil {
	}
	writer := io.MultiWriter(bl, os.Stdout)

	messageCh := make(chan []byte, 100)
	go handleMessages(writer, messageCh)

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGINT, syscall.SIGHUP)

	addr, _ := net.ResolveUDPAddr("udp4", ":7442")
	udpListener, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error starting UDP listener:", err.Error())
		return
	}

	go func() {
		var buffer []byte
		for {
			buffer = make([]byte, bufferSize)
			n, addr, err := udpListener.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println("Could not read from UDP conn:", err.Error())
				continue
			}

			fmt.Printf("received from %s: %s", addr.String(), buffer)
			messageCh <- buffer[:n]
		}
	}()

	tcpListener, err := net.Listen("tcp", ":7441")
	if err != nil {
		fmt.Println("Error starting TCP listener:", err.Error())
		return
	}

	go func() {
		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				fmt.Println("Could not read from UDP conn:", err.Error())
				continue
			}

			go handleRequest(messageCh, conn)
		}
	}()

	fmt.Println("Started. Waiting for interrupt...")
	<-exitCh
	if err := udpListener.Close(); err != nil {
		fmt.Println("could not close UDP listener:", err.Error())
	}
	if err := tcpListener.Close(); err != nil {
		fmt.Println("could not close TCP listener:", err.Error())
	}
	fmt.Println("Closed all connections, exiting.")
}

func handleRequest(messageCh chan []byte, conn net.Conn) {
	buf := make([]byte, bufferSize)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Could not read from TCP conn:", err.Error())
		return
	}
	fmt.Printf("Received from %s: %s\n", conn.RemoteAddr().String(), buf)
	messageCh <- buf[:n]

	if err = conn.Close(); err != nil {
		fmt.Println("Could not close TCP connection:", err.Error())
	}
}

func handleMessages(w io.Writer, messageCh chan []byte) {
	for {
		select {
		case msg, ok := <-messageCh:
			if ok {
				w.Write(msg)
			}
		}
	}
}

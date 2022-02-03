package main

import (
	"context"
	"flag"
	"fmt"
	logRotator "github.com/KaiserWerk/go-log-rotator"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const bufferSize = 1024

var (
	Version = "DEV"
	Commit  = "00000000"
	Date    = "0000-00-00T00:00:00.000"

	logPath = flag.String("logpath", ".", "The directory to place the collected logs into.")
	udpPort = flag.Int("udp", 7442, "The UDP Port to use")
	tcpPort = flag.Int("tcp", 7441, "The TCP Port to use")
)

func main() {
	flag.Parse()

	w, err := logRotator.New(*logPath, "loggy.log", 10<<20, 0644, 0, true)
	if err != nil {
		fmt.Println("could not create log rotator:", err.Error())
		return
	}

	messageCh := make(chan []byte, 100)
	go handleMessages(w, messageCh)

	ctx, cancel := context.WithCancel(context.Background())

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGINT, syscall.SIGHUP)

	fmt.Println("Starting up Loggly")
	fmt.Println("\tVersion:", Version)
	fmt.Println("\tCommit: ", Commit)
	fmt.Println("\tDate:   ", Date)

	fmt.Printf("Starting UDP Listener on port %d...\n", *udpPort)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", *udpPort))
	udpListener, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error starting UDP listener:", err.Error())
		return
	}

	go func() {
		var buffer []byte
		for {
			buffer = make([]byte, bufferSize)
			select {
			case <-ctx.Done():
				return
			default:
				n, addr, err := udpListener.ReadFromUDP(buffer)

				if err != nil {
					fmt.Println("Could not read from UDP conn:", err.Error())
					continue
				}

				fmt.Printf("received from %s: %s\n", addr.String(), buffer)
				messageCh <- buffer[:n]
			}

		}
	}()

	fmt.Printf("Starting TCP Listener on port %d...\n", *tcpPort)
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", *tcpPort))
	if err != nil {
		fmt.Println("Error starting TCP listener:", err.Error())
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				conn, err := tcpListener.Accept()
				if err != nil {
					fmt.Println("Could not read from UDP conn:", err.Error())
					continue
				}

				go handleRequest(messageCh, conn)
			}
		}
	}()

	fmt.Println("Started. Waiting for interrupt...")
	<-exitCh
	cancel()

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
				_, err := fmt.Fprintf(w, "%s %s\n", time.Now().Format(time.RFC3339), strings.TrimSpace(string(msg)))
				if err != nil {
					fmt.Println("could not execute write:", err.Error())
				}
			}
		}
	}
}

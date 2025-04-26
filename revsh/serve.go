package revsh

import (
	"log"
	"net"
	"time"
)

func ServeTCP(address, cmd string, timeout int) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		} else {
			log.Printf("TCP connected %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		}
		go handleTCP(conn, cmd, timeout)
	}
}

func handleTCP(conn net.Conn, cmd string, timeout int) {
	defer conn.Close()

	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	conn.SetDeadline(deadline)

	_, err := conn.Write([]byte(cmd + "\n"))
	if err != nil {
		log.Printf("Write error: %v", err)
		return
	}
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Println("Connection timed out")
			} else {
				log.Printf("Read error: %v\n", err)
			}
			return
		}
		log.Printf("Received: %d %s\n", n, LogEscape(string(buf[:n])))
	}
}

func ServeUDP(address, cmd string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	var conn *net.UDPConn
	conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	cmd = cmd + "\n"

	var buf [2048]byte
	for {
		rlen, remote, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			return err
		} else {
			log.Printf("UDP Received %s - %d %s\n", remote, rlen, LogEscape(string(buf[:rlen])))
			_, err = conn.WriteToUDP([]byte(cmd), remote)
			if err != nil {
				log.Printf("failed to write udp data: %v\n", err)
			}

		}
	}
}

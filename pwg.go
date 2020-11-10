package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
)

const (
	connSockPort = ":8712"
	connSockType = "tcp"
	connWsPort   = ":8715"
)

var (
	wsg         *websocket.Conn
	allTCPConns []tcpConn
)

type tcpConn struct {
	clientConn net.Conn
	clientAddr net.Addr
}

func handleConnection(conn net.Conn, k int) {
	var bufTrimmed []byte
	for {
		buf := make([]byte, 1024)
		k := 0
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			conn.Close()
			break
		}

		bufTrimmed = rmvNull(buf)
		if wsg != nil {
			websocket.Message.Send(wsg, fmt.Sprint(rmvNull(bufTrimmed)))
		}
		fmt.Println(rmvNull(bufTrimmed))
		fmt.Println(fmt.Sprint(k))
		if string(bufTrimmed) == "GetID" {
			conn.Write([]byte(fmt.Sprintf("ID:%d", k)))
		}
		k++
	}

}

func rmvNull(toRmv []byte) []byte {
	var tmpByte []byte
	for _, b := range toRmv {
		if b != 0 {
			tmpByte = append(tmpByte, b)
		}
	}
	return tmpByte
}

func main() {
	i := 0
	var tmpClient tcpConn
	ln, err := net.Listen(connSockType, connSockPort)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Accept connections on port %s for tcp and %s for websocket", connSockPort, connWsPort))
	http.Handle("/", websocket.Handler(websock))

	go http.ListenAndServe(connWsPort, nil)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		tmpClient.clientAddr = conn.RemoteAddr()
		tmpClient.clientConn = conn
		allTCPConns = append(allTCPConns, tmpClient)
		k := i
		i++
		go handleConnection(allTCPConns[k].clientConn, k)

	}
}

func websock(ws *websocket.Conn) {
	wsg = ws
	var err error
	var reply string
	for {
		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)
		sendToTCP(reply)
		msg := "Received: " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(wsg, msg); err != nil {
			fmt.Println("Can't send")
			break
		}

	}
}

func sendToTCP(msg string) {
	tmpStrSlice := strings.Split(msg, ">")
	id, _ := strconv.ParseInt(tmpStrSlice[0], 0, 64)
	smsg := tmpStrSlice[1]
	allTCPConns[id].clientConn.Write([]byte(smsg))

	/*for i, tcp := range allTCPConns {
		tcp.clientConn.Write([]byte(msg + fmt.Sprint(i)))
	}*/
}

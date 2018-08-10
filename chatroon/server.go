package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

var onlineConn = make(map[string]net.Conn)
var messageQueue = make(chan string, 1000)
var quitChann = make(chan bool)

//LOG_DIRECTORY

func checkError(err error) {
	if err != nil {
		fmt.Printf("err:%s", err.Error())
		os.Exit(1)
	}
}

func consumeMessage() {
	for {
		select {
		case message := <-messageQueue:
			doProcessMessage(message)
		case <-quitChann:
			break
		}
	}
}

func doProcessMessage(message string) {
	contents := strings.Split(message, "#")
	if len(contents) > 1 {
		addr := contents[0]
		messageContent := strings.Join(contents[1:], "#")
		if conns, ok := onlineConn[addr]; ok {
			_, err := conns.Write([]byte(messageContent))
			checkError(err)
		}
	} else {
		//IP:PORT*LIST
		contentss := strings.Split(message, "*")
		if strings.ToUpper(contentss[1]) == "LIST" {
			var ips string = ""
			for k, _ := range onlineConn {
				ips = ips + "|" + k
			}
			if conn, ok := onlineConn[contentss[0]]; ok {
				_, err := conn.Write([]byte(ips))
				checkError(err)
			}
		}
	}
}

func messageReceive(conn net.Conn) {
	buf := make([]byte, 1024)
	defer func(conn net.Conn) {
		delete(onlineConn, conn.RemoteAddr().String())
		conn.Close()
		fmt.Println(conn.RemoteAddr().String() + "已经离开")
		for k, _ := range onlineConn {
			fmt.Print("目前的成员：" + k + " ")
		}
		fmt.Println()
	}(conn)

	for {
		content, err := conn.Read(buf)
		checkError(err)
		if content != 0 {
			messageQueue <- string(buf[:content])
		}
	}
}

func main() {
	listen_socket, err := net.Listen("tcp", "127.0.0.1:8081")
	checkError(err)
	defer listen_socket.Close()
	fmt.Println("Server Started...")
	go consumeMessage()
	for {
		conn, err := listen_socket.Accept()
		checkError(err)
		onlineConn[conn.RemoteAddr().String()] = conn
		for k := range onlineConn {
			fmt.Print("目前的成员：" + k + " ")
		}
		fmt.Println()
		go messageReceive(conn)
	}
}

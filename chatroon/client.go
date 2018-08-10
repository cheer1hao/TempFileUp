package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func checkErrorS(err error) {
	if err != nil {
		fmt.Printf("err:%s", err.Error())
		os.Exit(1)
	}
}

func sendMessage(conn net.Conn) {
	var input string
	defer conn.Close()
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)
		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		}
		_, err := conn.Write([]byte(input))
		checkErrorS(err)
	}
}

func main() {
	listen_socket, err := net.Dial("tcp", "127.0.0.1:8081")
	checkErrorS(err)
	go sendMessage(listen_socket)
	buf := make([]byte, 1024)
	for {
		length, err := listen_socket.Read(buf)
		//用户退出后此处会报错，因此要平滑处理
		if err != nil {
			fmt.Println(listen_socket.RemoteAddr().String() + ":你已退出")
			os.Exit(0)
		}
		fmt.Println("receive server message:" + string(buf[:length]))
	}
}

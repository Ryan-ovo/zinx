package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("client start...")
	time.Sleep(time.Second)
	// 1. 直接连接远程服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for {
		// 2. 调用write写数据
		_, err = conn.Write([]byte("Hello zinx v0.1"))
		if err != nil {
			fmt.Println("write conn error")
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf error")
			return
		}
		fmt.Printf("server call back: %s, cnt: %d\n", buf, cnt)
		time.Sleep(time.Second)
	}

}

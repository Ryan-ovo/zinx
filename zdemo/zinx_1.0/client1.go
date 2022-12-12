package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"zinx/znet"
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
		// 发送封包之后的消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(2, []byte("Zinx V1.0 Client1 Test Message")))
		if err != nil {
			log.Println("pack error = ", err)
			return
		}
		if _, err = conn.Write(binaryMsg); err != nil {
			log.Println("write error = ", err)
			return
		}

		// 解析服务器的返回
		msgHead := make([]byte, dp.GetHeadLen())
		if _, err = io.ReadFull(conn, msgHead); err != nil {
			log.Println("read message head error = ", err)
			return
		}
		msg, err := dp.Unpack(msgHead)
		if err != nil {
			log.Println("unpack msg from server error = ", err)
			return
		}
		if msg.GetMsgLen() > 0 {
			msgBody := make([]byte, msg.GetMsgLen())
			if _, err = io.ReadFull(conn, msgBody); err != nil {
				log.Println("read msg body error = ", err)
				return
			}
			msg.SetData(msgBody)
			fmt.Println("---> Receive Server Msg: ID = ", msg.GetMsgID(), ", len = ", msg.GetMsgLen(), ", data = ", string(msg.GetData()))
		}

		time.Sleep(time.Second)
	}
}

package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 注释GlobalObject.Reload()
func TestDataPack(t *testing.T) {
	// 1. 模拟服务端
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Error(err)
	}
	go func() {
		for {
			// 循环处理客户端连接请求
			conn, err := listener.Accept()
			if err != nil {
				t.Error(err)
			}
			go func(conn net.Conn) {
				// 创建封包拆包处理对象
				dp := NewDataPack()
				for {
					// 读取分两次，先读头部，再读消息体
					headData := make([]byte, dp.GetHeadLen())
					_, err = io.ReadFull(conn, headData)
					if err != nil {
						t.Error(err)
					}
					// 先把头部信息拆包到Message对象
					msg, err := dp.Unpack(headData)
					if err != nil {
						t.Error(err)
					}
					// 第二次读取消息体
					msgBody := make([]byte, msg.GetMsgLen())
					_, err = io.ReadFull(conn, msgBody)
					if err != nil {
						t.Error(err)
					}
					msg.SetData(msgBody)
					// 本次拆包成功
					t.Log("---> Receive MsgID: ", msg.GetMsgID(), ", DataLen = ", msg.GetMsgLen(), "Data = ", string(msg.GetData()))
				}
			}(conn)
		}
	}()
	// 2. 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}
	// 创建一个封包对象 dp
	dp := NewDataPack()

	// 模拟粘包过程，封装两个msg一起发送
	// 封装msg1
	msg1 := &Message{
		ID:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error: ", err)
		return
	}
	// 封装msg2
	msg2 := &Message{
		ID:      2,
		DataLen: 6,
		Data:    []byte{'h', 'e', 'l', 'l', 'o', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		t.Error("client pack msg2 error: ", err)
		return
	}
	// 将两个包粘在一起
	sendData1 = append(sendData1, sendData2...)
	// 一次性发送给服务端
	_, err = conn.Write(sendData1)
	if err != nil {
		t.Error(err)
	}
	// 客户端阻塞
	select {}
}

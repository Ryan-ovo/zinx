package znet

import (
	"errors"
	"log"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	c.connections[conn.GetConnID()] = conn
	log.Printf("Add connection success! connID = [%d], connNum = [%d]\n", conn.GetConnID(), c.Len())
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	delete(c.connections, conn.GetConnID())
	log.Printf("Remove connection success! connID = [%d], connNum = [%d]\n", conn.GetConnID(), c.Len())
}

func (c *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}
	log.Printf("Clear all connections successfully: conn num = [%d]\n", c.Len())
}

func (c *ConnManager) ClearOneConn(connID uint32) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	if conn, ok := c.connections[connID]; ok {
		//停止
		conn.Stop()
		//删除
		delete(c.connections, connID)
		log.Println("Clear Connections ID:  ", connID, "successfully")
		return
	}
	log.Println("Clear Connections ID:  ", connID, "failed")
}

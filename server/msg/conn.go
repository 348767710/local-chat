package msg

import (
	"github.com/gorilla/websocket"
	"sync"
	"errors"
	"log"
)

// 用户结构体
type User struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Status int `json:"status"`
	Sign string `json:"sign"`
	Avatar string `json:"avatar"`
}

// 群组结构体
type Group struct {
	Id int `json:"id"`
	Groupname string `json:"groupname"`
	Avatar string `json:"avatar"`
}

// 用户连接
type Conn struct {
	sync.Mutex
	Conn      *websocket.Conn
	writeChan chan []byte
	maxMsgLen uint32
	closeFlag bool
}

func (c *Conn) WriteMsg(args ...[]byte) error {
	c.Lock()
	defer c.Unlock()
	if c.closeFlag {
		return nil
	}

	// get len
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}

	// check len
	if msgLen > c.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < 1 {
		return errors.New("message too short")
	}

	// don't copy
	if len(args) == 1 {
		c.doWrite(args[0])
		return nil
	}

	// merge the args
	msg := make([]byte, msgLen)
	l := 0
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}

	c.doWrite(msg)

	return nil
}

func (c *Conn) doWrite(b []byte) {
	if len(c.writeChan) == cap(c.writeChan) {
		log.Print("close conn: channel full")
		return
	}

	c.writeChan <- b
}

// goroutine not safe
func (c *Conn) ReadMsg() ([]byte, error) {
	_, b, err := c.Conn.ReadMessage()
	log.Println("1. 收到消息: --> " ,string(b))
	return b, err
}


func NewConn(conn *websocket.Conn) *Conn {
	c := new(Conn)
	c.Conn = conn
	c.writeChan = make(chan []byte,10240)
	c.maxMsgLen = 40960
	go func() {
		for b := range c.writeChan {
			if b == nil {
				break
			}
			err := conn.WriteMessage(websocket.TextMessage,b)
			if err != nil {
				break
			}
		}
		conn.Close()
		c.Lock()
		c.closeFlag = true
		c.Unlock()
	}()

	return c
}
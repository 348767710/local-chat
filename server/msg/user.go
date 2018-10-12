package msg

import (
	"sync"
	"sync/atomic"
	"encoding/json"
	"log"
)

var (
	Conns = make(map[int64]*Conn,0)
	ConnMutex = sync.RWMutex{}
	Id int64
	IdUserMap = make(map[string]int64,0)
	IdUserMapMutex = sync.RWMutex{}
)

func GetNewId() int64 {
	return atomic.AddInt64(&Id,1)
}

// 新增连接
func AddNewConn(id int64,username string,conn *Conn)  {

	if id , ok := IdUserMap[username] ; ok {
		if conn , ok := Conns[id] ; ok {
			conn.Lock()
			conn.closeFlag = true
			conn.Unlock()
			conn.Conn.Close() // 关闭老的连接
		}
	}

	ConnMutex.RLock()
	Conns[id] = conn
	ConnMutex.RUnlock()

	IdUserMapMutex.RLock()
	IdUserMap[username] = id
	IdUserMapMutex.RUnlock()
}

func SendToId(id int64,v interface{}) (error)  {
	b , err := json.Marshal(v)
	if err != nil {
		return err
	}
	return Conns[id].WriteMsg(b)
}
// 群发给所有人
func SendToAllExceptSelf(id int64, v interface{}) (error) {
	b , err := json.Marshal(v)
	if err != nil {
		return err
	}
	for k, val := range Conns {
		if k == id {
			continue
		}
		log.Println("6. 发送给id ->",k)
		err = val.WriteMsg(b)
		if err != nil {
			return err
		}
	}
	return nil
}
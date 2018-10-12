package msg

import (
	"fmt"
	"log"
)

type Msg struct {
	MsgType int `json:"msg_type"`// 消息类型
	Body interface{} `json:"body"`// 消息体.
}

type IMsg interface{
	Send(conn *Conn) error
}

type ConnectMsg struct{
	Username string `json:"username"`
	Id int64 `json:"id"`
}
// 新用户连接.
func (cm *ConnectMsg)Send(conn *Conn) error {
	cm.Id = GetNewId()
	AddNewConn(cm.Id,cm.Username,conn)

	var msg = new(Msg)
	msg.MsgType = Connect
	msg.Body = cm

	// 发给自己
	err := SendToId(cm.Id,msg)
	if err != nil {
		return err
	}
	log.Println("3.发送给了自己,内容是：",msg)
	// 群发.
	sys := &System{
		System:true,
		Id:1 ,
		Type: "group",
		Content: fmt.Sprintf("%s加入群聊",cm.Username),
	}
	msg.MsgType = SYSTEM
	msg.Body = sys
	log.Println("4.发送系统消息：内容：",msg)
	return SendToAllExceptSelf(cm.Id,msg)
}

type GroupMsg struct{
	Id int64 `json:"id"`
	Username string `json:"username"`
	Avatar string `json:"avatar"`
	Type string `json:"type"`
	Content string `json:"content"`
	SelfId int64 `json:"self_id"`
}

func (gm *GroupMsg)Send(conn *Conn) error {
	var msg = new(Msg)
	msg.MsgType = GroupMessage
	msg.Body = gm
	log.Println("5.发送群组消息：内容：",msg)
	return SendToAllExceptSelf(gm.SelfId,msg)
}

type System struct {
	System bool `json:"system"`
	Id int `json:"id"`
	Type string `json:"type"`
	Content string `json:"content"`
}


func (m *Msg) Parse () (IMsg,error) {

	return nil,nil
}

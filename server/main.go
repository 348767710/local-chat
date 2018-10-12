package main

import (
	"flag"
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"local-chat/server/msg"
	"encoding/json"
)

var (
	staticDir = flag.String("static","","staic文件夹")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func main()  {
	flag.Parse()

	static := http.FileServer(http.Dir(*staticDir))
	log.Println(*staticDir)

	server := http.DefaultServeMux
	server.HandleFunc("/ws",handler)
	server.Handle("/", static)
	server.HandleFunc("/getList.json", getOnlineList)

	go printConn()

	if err := http.ListenAndServe(":7111",server); err != nil {
		log.Fatal(err)
	}

}

func handler(w http.ResponseWriter, r *http.Request)  {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	mConn := msg.NewConn(conn)

	for {
		data , err := mConn.ReadMsg()

		if err != nil {
			log.Println("read error 1:",err)
			break
		}
		var m = new(msg.Msg)
		err = json.Unmarshal(data,m)
		if err != nil {
			log.Println(err)
		}
		var im msg.IMsg

		switch m.MsgType {
		case msg.Connect:
			im = new(msg.ConnectMsg)
		case msg.GroupMessage:
			im = new(msg.GroupMsg)
		}
		log.Println("2. 消息类型是：-->",m.MsgType ," 内容：body--->",m.Body)
		b ,err := json.Marshal(m.Body)
		if err != nil {
			break
		}
		err = json.Unmarshal(b,im)
		if err != nil {
			break
		}
		err = im.Send(mConn)
		if err != nil {
			break
		}
	}
}

// getOnlineList 获取在线列表
func getOnlineList(w http.ResponseWriter, r *http.Request)  {

}

func printConn()  {

}
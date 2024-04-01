package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId   uint   // 发送者id
	TargetId uint   // 接收者id
	Type     uint   // 发送类型 1私聊 2群聊 3广播
	Media    uint   // 消息类型 1文字 2图片 3音频
	Content  string // 内容
	Pic      string // 图片
	Url      string // 链接
	Amount   int    //
}

func (msg *Message) TableName() string {
	return "message"

}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// 发送者Id 接受者Id 发送类型 消息类型 消息内容
func Chat(writer http.ResponseWriter, request *http.Request) {

	// 1. 获取参数
	query := request.URL.Query()

	userId, _ := strconv.ParseInt(query.Get("userId"), 10, 64)
	targetId, _ := strconv.ParseInt(query.Get("targetId"), 10, 64)
	// context := query.Get("context")
	// sendType := query.Get("type")

	isValid := true
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		fmt.Println(err)

		return
	}

	// 2. 获取 conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 3. 用户关系

	// 4. userId 与 node绑定 并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 5. 完成发送逻辑
	go sendProc(node)

	// 6. 完成接收逻辑
	go receiveProc(node)

	sendMsg(uint(userId), uint(targetId), []byte("欢迎进入聊天室"))

}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func receiveProc(node *Node) {
	for {

		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		broadMsg(data)

		// 处理数据
		fmt.Println("[ws]<<<<<<<", data)

	}

}

var udpsendChan chan []byte = make(chan []byte, 10)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpReceiveProc()
}

func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 1, 1),
		Port: 3000,
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		select {
		case data := <-udpsendChan:
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
func udpReceiveProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println(err)
			return
		}
		disPatch(buf[:n])
	}

}
func disPatch(data []byte) {
	msg := Message{}
	// json 转换
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch msg.Type {
	case 1: // 私聊
		sendMsg(msg.FromId, msg.TargetId, data)
		// case 2:
		// 	sendGroupMsg(msg)
		// case 3:
		// 	broadMsg(data)
	}

}

func sendMsg(userId uint, targetId uint, msg []byte) {
	// 1. 获取接收者的node
	rwLocker.RLock()
	node, ok := clientMap[int64(userId)]
	rwLocker.RUnlock()
	if !ok {
		fmt.Println("用户不在线")
		return
	}
	node.DataQueue <- msg
}

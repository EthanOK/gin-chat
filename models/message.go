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
	UserId   uint   // 发送者id
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
	isValid := true

	// 2. 建立一个 WebSocket 连接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		fmt.Println(err)

		return
	}

	// 3. 定义 WebSocket 节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 4. userId 与 node绑定 并加锁 谁进来谁在线
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 5. 完成发送消息到 WebSocket 逻辑
	go sendProc(node)

	// 6. 完成从 WebSocket 接收消息逻辑
	go receiveProc(node)

	// 7. 后台向登录用户发送欢迎消息
	hello := "欢迎用户" + FindNameByUserId(uint(userId)) + "进入聊天室"
	sendMsg(uint(userId), []byte(hello))

}

func sendProc(node *Node) {
	// 从 node.DataQueue 通道中读取数据，并将其发送到 WebSocket 中
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws] sendProc >>>>", string(data))
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
		// 从 Websocket 中读取一条消息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		// 处理接受到的数据
		disPatch(data)
		// 广播数据
		// broadMsg(data)

		fmt.Println("[ws] receiveProc <<<<", string(data))

	}
}

var udpsendChan chan []byte = make(chan []byte, 10)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpReceiveProc()
	fmt.Println("init success!!!")
}

// udp 数据发送
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 255),
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
			fmt.Println("[udp] sendProc >>>>", string(data))
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

		fmt.Println("[udp] receiveProc <<<<", string(buf[:n]))

		disPatch(buf[:n])
	}

}

// 后端调度逻辑处理
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
		sendMsg(msg.TargetId, data)
		// case 2:
		// 	sendGroupMsg(msg)
		// case 3:
		// 	broadMsg(data)
	}
}

func sendMsg(userId uint, msg []byte) {
	fmt.Println("sendMsg >>> userId:", userId, "message:", string(msg))

	rwLocker.RLock()
	node, ok := clientMap[int64(userId)]
	rwLocker.RUnlock()

	if ok {
		node.DataQueue <- msg
	}
}

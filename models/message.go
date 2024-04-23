package models

import (
	"context"
	"encoding/json"
	"fmt"
	"gin-chat/utils"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
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
	Addr      string //客户端地址
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
	userId_ := query.Get("userId")
	userId, _ := strconv.ParseInt(userId_, 10, 64)
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
		Addr:      conn.RemoteAddr().String(), //客户端地址
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

	//7.加入在线用户到缓存
	SetUserOnlineInfo("online_"+userId_, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)

	// // 8. 后台向登录用户发送欢迎消息
	// hello := "欢迎用户" + FindNameByUserId(uint(userId)) + "进入聊天室"
	// sendMsg(uint(userId), []byte(hello))

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
		// 广播消息到局域网
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

	jsonMsg := Message{}
	json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.UserId))
	r, err := utils.Redis.Get(ctx, "online_"+userIdStr).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r)

	if r != "" {
		if ok {
			node.DataQueue <- msg
		}
	}

	var key string
	if userId > jsonMsg.UserId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}

	// 从 Redis 中以逆序方式获取指定键对应的有序集合中的所有成员
	res, err := utils.Redis.ZRevRange(ctx, key, 0, -1).Result()

	if err != nil {
		fmt.Println(err)
	}

	// 往 Redis 中指定键对应的有序集合中添加一条新消息，新消息的分数比已有成员的最高分数高，确保它被放在最前面
	score := float64(cap(res)) + 1
	ress, e := utils.Redis.ZAdd(ctx, key, redis.Z{Score: score, Member: msg}).Result()
	//res, e := utils.Red.Do(ctx, "zadd", key, 1, jsonMsg).Result() //备用 后续拓展 记录完整msg
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(ress)
}

// 获取缓存里面的消息
func RedisMsg(userIdA int64, userIdB int64, start int64, end int64, isRev bool) []string {
	rwLocker.RLock()
	//node, ok := clientMap[userIdA]
	rwLocker.RUnlock()
	//jsonMsg := Message{}
	//json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userIdA))
	targetIdStr := strconv.Itoa(int(userIdB))
	var key string
	if userIdA > userIdB {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}
	//key = "msg_" + userIdStr + "_" + targetIdStr
	//rels, err := utils.Red.ZRevRange(ctx, key, 0, 10).Result()  //根据score倒叙

	var rels []string
	var err error
	if isRev {
		rels, err = utils.Redis.ZRange(ctx, key, start, end).Result()
	} else {
		rels, err = utils.Redis.ZRevRange(ctx, key, start, end).Result()
	}
	if err != nil {
		fmt.Println(err) //没有找到
	}
	// 发送推送消息
	/**
	// 后台通过websoket 推送消息
	for _, val := range rels {
		fmt.Println("sendMsg >>> userID: ", userIdA, "  msg:", val)
		node.DataQueue <- []byte(val)
	}**/
	return rels
}

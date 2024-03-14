package utils

import (
	"context"
	"fmt"
)

const (
	PublishChannel = "websocket"
)

// 发布消息到 Redis
func Publish(ctx context.Context, channel string, message string) error {

	fmt.Println("Publish: ")

	return Redis.Publish(ctx, channel, message).Err()

}

// 订阅 Redis 消息
func Subscribe(ctx context.Context, channel string) (string, error) {

	sub := Redis.Subscribe(ctx, channel)
	fmt.Println("Subscribe0: ", ctx)

	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		return "", err
	}

	fmt.Println("Subscribe1: ", msg.String())

	return msg.String(), err
}

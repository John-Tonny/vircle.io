package main

import (
	"encoding/json"

	"github.com/micro/go-micro/broker"
	nodePB "vircle.io/mnhost/interface/out/nodemsg"

	//"context"
	"fmt"
	"log"

	//"time"

	"github.com/micro/go-micro"

	//"github.com/micro/go-micro/util/log"

	"vircle.io/mnhost/vps/handler"
	"vircle.io/mnhost/vps/subscriber"

	//aaa "vircle.io/mnhost/vps/handler"
	vps "vircle.io/mnhost/interface/proto/vps"

	"vircle.io/mnhost/common"
	//"vircle.io/mnhost/config"
)

const serviceName = "vps"

func main() {
	fmt.Println("start vps")
	session, err := common.CreateDBSession(serviceName)
	fmt.Println("session:", session)

	if err != nil {
		log.Fatalf("create session error: %v\n", err)
	}

	// 创建于 Redis 的主会话，需在退出 main() 时候手动释放连接
	defer session.Close()

	srv := common.GetMicroServer(serviceName)

	bk := srv.Server().Options().Broker

	// 将实现服务端的 API 注册到服务端
	vps.RegisterVpsHandler(srv.Server(), handler.GetHandler(session, bk))

	// Register Handler
	//vps.RegisterVpsHandler(srv.Server(), new(handler.Vps))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("Vircle.Mnhost.vps", srv.Server(), new(subscriber.Vps))
	/*
		// Register Function as Subscriber
		micro.RegisterSubscriber("go.micro.srv.vps1a", srv.Server(), subscriber.Handler)

		// Register Function as Subscriber
		micro.RegisterSubscriber("go.micro.srv.vps.newnode", srv.Server(), subscriber.NewNode)

		// Register Function as Subscriber
		micro.RegisterSubscriber("go.micro.srv.vps.delnode", srv.Server(), subscriber.DelNode)
	*/
	//pubMsg("1", "Vircle.Mnhost.Topic.NewOrder", 1, bk)

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func pubMsg(userID, topic string, msgId int64, e broker.Broker) error {
	fmt.Println("start msg publish")
	msg := nodePB.NodeMsg{
		Id: msgId,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	data := &broker.Message{
		Header: map[string]string{
			"user_id": userID,
		},
		Body: body,
	}

	if err := e.Publish(topic, data); err != nil {
		log.Printf("[pub] failed: %v\n", err)
	}
	fmt.Println("success msg publish:", msgId)
	return nil
}

package main

import (
	"fmt"
	//"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	//"github.com/micro/go-plugins/broker/redis"

	"vircle.io/mnhost/vps/handler"
	"vircle.io/mnhost/vps/subscriber"

	//aaa "vircle.io/mnhost/vps/handler"
	vps "vircle.io/mnhost/interface/out/vps"

	"github.com/micro/go-micro/broker"
	//"vircle.io/mnhost/go-plugins/broker/kafka"
	//"github.com/micro/go-micro/registry"
	//"github.com/micro/go-plugins/registry/consul"
)

func main() {
	/*b := kafka.NewBroker(
		broker.Addrs("127.0.0.1:9092"),
	)

		b := nsq.NewBroker(
			broker.Addrs("redis://localhost:6379/"),
		)

	b.Init()
	b.Connect()*/
	/*
		msg := &broker.Message{
			Header: map[string]string{
				"id": fmt.Sprintf("%d", 10),
			},
			Body: []byte(fmt.Sprintf("%d: %s", 10, time.Now().String())),
		}
		err := b.Publish("topic string test", msg)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("publish")
		log.Info(b)*/

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.vps"),
		micro.Version("latest"),
		/*micro.Broker(kafka.NewBroker(func(o *broker.Options) {
			o.Addrs = []string{"127.0.0.1:9092"}
		})),*/
		//micro.Broker(b),
	)

	if err := broker.Connect(); err != nil {
		fmt.Println(err.Error())
	}

	_, err := broker.Subscribe("my test", func(p broker.Event) error {
		fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	// Initialise service
	service.Init()

	// Register Handler
	vps.RegisterVpsHandler(service.Server(), new(handler.Vps))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.vps", service.Server(), new(subscriber.Vps))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.vps", service.Server(), subscriber.Handler)

	/*
		_, err = b.Subscribe("topic strint test", func(p broker.Event) error {
			fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
			return nil
		})
		if err != nil {
			fmt.Println("error")
			fmt.Println(err)
		}*/

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

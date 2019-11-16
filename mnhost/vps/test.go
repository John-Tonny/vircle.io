package main

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/util/log"

	//"github.com/micro/go-plugins/broker/redis"

	"vircle.io/mnhost/vps/handler"
	"vircle.io/mnhost/vps/subscriber"

	vps "vircle.io/mnhost/interface/out/vps"
	aaa "vircle.io/mnhost/vps/handler"

	"github.com/micro/go-micro/broker"
	//"vircle.io/mnhost/go-plugins/broker/kafka"

	//"github.com/micro/go-micro/registry"
	//"github.com/micro/go-plugins/registry/consul"
	"github.com/micro/go-micro/server"
)

var (
	topic = "mytest"
)

func main() {

	var vpsInfo aaa.VpsInfo

	aaa.UpdateVps("", "", "", 1, vpsInfo)
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

	//aaa.WriteConf("/tmp/vircle.conf", 8)
	//aaa.NewNode("", "192.168.246.180", "htjonny")

	//vpsInfo, err := aaa.NewVps("", "", "")
	//log.Info(vpsInfo)
	//log.Info(err)
	//aaa.Testuser("eee", "eee", "13510019994")
	/*aaa.Testorder("13510019994")
	aaa.Testuser("fff", "fff", "13510019993")
	aaa.Testorder("13510019993")
	aaa.TestProduct()
	aaa.TestCoin()*/
	//aaa.Testnode("2")

	//aaa.NewNode(5, "vircle", "18.222.181.64", "vpub$999000")

	// Initialise service
	service.Init()

	if err := broker.Connect(); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("start service")

	// Register Handler
	vps.RegisterVpsHandler(service.Server(), new(handler.Vps))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.vps2a", service.Server(), new(subscriber.Vps))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.vps3a", service.Server(), subscriber.Handler)

	p := micro.NewPublisher("go.micro.srv.vpsa", service.Client())
	pub2("mytest1", p)

	p1 := micro.NewPublisher("go.micro.srv.vps.newnode", service.Client())
	pub2("mytest2", p1)

	p2 := micro.NewPublisher("go.micro.srv.vps.newnode", service.Client())
	pub2("mytest3", p2)

	sub1()

	/*
		_, err = b.Subscribe("topic strint test", func(p broker.Event) error {
			fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
			return nil
		})
		if err != nil {
			fmt.Println("error")
			fmt.Println(err)
		}*/

	//go pub()
	go sub()
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func pub() {
	tick := time.NewTicker(time.Second)
	i := 0
	for _ = range tick.C {
		msg := &broker.Message{
			Header: map[string]string{
				"id": fmt.Sprintf("%d", i),
			},
			Body: []byte(fmt.Sprintf("%d: %s ---%s", i, time.Now().String(), "lht")),
		}
		if err := broker.Publish(topic, msg); err != nil {
			fmt.Println("[pub] failed: %v", err)
		} else {
			fmt.Println("[pub] pubbed message:", string(msg.Body))
		}
		i++
	}
}

func sub() {
	_, err := broker.Subscribe(topic, func(p broker.Event) error {
		fmt.Println("[sub] received message:", string(p.Message().Body), "header", p.Message().Header)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func sub1() {
	// Register Subscribers
	if err := server.Subscribe(
		server.NewSubscriber(
			"topic.examplea",
			new(subscriber.Vps),
		),
	); err != nil {
		log.Fatal(err)
	}
}

func pub2(topic string, p micro.Publisher) {
	ev := &vps.Message{
		Say: "aaattt:" + topic,
	}

	fmt.Println("publishing %+v\n", ev)

	// publish an event
	if err := p.Publish(context.Background(), ev); err != nil {
		log.Logf("error publishing: %v", err)
	}
}

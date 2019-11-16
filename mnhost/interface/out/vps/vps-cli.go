package main

import (
	"fmt"
	"log"

	"context"
	//"reflect"

	"vircle.io/mnhost/common"
	"vircle.io/mnhost/config"
	pb "vircle.io/mnhost/interface/out/vps"
)

const service = "vps"

var (
	serviceName string
)

func init() {
	serviceName = config.GetServiceName(service)
}

func main() {
	srv := common.GetMicroClient(service)

	//fmt.Println(reflect.TypeOf(srv))
	fmt.Println(service, ":", serviceName)

	// 创建 user-service 微服务的客户端
	client := pb.NewVpsService(serviceName, srv.Client())

	//name := "Ethan"

	resp, err := client.Call(context.Background(), &pb.Request{
		Id: 1,
	})
	if err != nil {
		log.Printf("call Create error: %v", err)
	} else {
		log.Println("created: ", "success call")
	}
	fmt.Println(resp)

	/*
		var id int64
		id = 3
		resp, err := client.NewNode(context.Background(), &pb.Request{
			Id: id,
		})
		if err != nil {
			log.Printf("call Create error: %v", err)
		} else {
			log.Println("created: ", "success call")
		}
		fmt.Println(resp)
	*/
}

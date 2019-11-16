package main

import (
	"fmt"
	//"log"

	//"vircle.io/mnhost/common"
	//"vircle.io/mnhost/config"

	pb "vircle.io/mnhost/interface/out/vps"
	//"context"
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
	fmt.Println(srv)
	fmt.Println(service, ":", serviceName)

	/*
		// 创建 user-service 微服务的客户端
		client := pb.NewVpService(serviceName, srv.Client())

		name := "Ethan"

		resp, err := client.Call(context.Background(), &pb.Message{
			Say: name,
		})
		if err != nil {
			log.Printf("call Create error: %v", err)
		} else {
			log.Println("created: ", "success call")
		}*/

}

package main

import (
	"log"

	"vircle.io/mnhost/common"
	"vircle.io/mnhost/config"
	pb "vircle.io/mnhost/interface/out/user"

	"context"
)

const service = "user"

var (
	serviceName string
)

func init() {
	serviceName = config.GetServiceName(service)
}

func main() {
	srv := common.GetMicroClient(service)

	// 创建 user-service 微服务的客户端
	client := pb.NewUserService(serviceName, srv.Client())

	name := "Ethan"
	password := "test123"
	mobile := "13588889997"
	realName := "test-company"
	idCard := "222233334444555566667777"

	resp, err := client.Create(context.Background(), &pb.User{
		Name:     name,
		Password: password,
		Mobile:   mobile,
		RealName: realName,
		IdCard:   idCard,
	})
	log.Println("start")
	if err != nil {
		log.Printf("call Create error: %v", err)
	} else {
		log.Println("created: ", resp.User.Id)
	}

	allResp, err := client.GetAll(context.Background(), &pb.Request{})
	if err != nil {
		log.Printf("call GetAll error: %v", err)
	} else {
		for i, u := range allResp.Users {
			log.Printf("user_%d: %v\n", i, u)
		}
	}

	authResp, err := client.Auth(context.Background(), &pb.User{
		Mobile:   mobile,
		Password: password,
	})
	if err != nil {
		log.Printf("auth failed: %v\n", err)
	} else {
		log.Println("token: ", authResp.Token)
	}
}

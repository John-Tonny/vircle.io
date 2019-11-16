package main

import (
	"fmt"
	"log"
	"time"

	"context"
	//"reflect"

	"vircle.io/mnhost/common"
	"vircle.io/mnhost/config"
	vps "vircle.io/mnhost/interface/proto/vps"
	//"github.com/astaxie/beego/orm"
	//"vircle.io/mnhost/model"
	//"github.com/dynport/gossh"
)

const service = "vps"

var (
	serviceName string
)

func init() {
	serviceName = config.GetServiceName(service)
}

func main() {
	fmt.Println("start")

	srv := common.GetMicroClient(service)

	//fmt.Println(reflect.TypeOf(srv))
	fmt.Println(service, ":", serviceName)

	/*
		client1 := gossh.New("18.224.46.202", "root")
		client1.SetPassword("vpub$999000")

		defer client1.Close()


			volumeName := "mn3"
			var rsp *gossh.Result
			cmd := "docker stop  `docker ps -aq --filter name=" + volumeName + "`"
			fmt.Println("cmd:", cmd)

			rsp, err := client1.Execute(cmd)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("rsp:", rsp.Stdout())

			var nodes models.Node
			o := orm.NewOrm()
			qs := o.QueryTable("node")
			err = qs.Filter("vps_id", 3).One(&nodes)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(nodes)
			fmt.Println("orderId:", nodes.OrderNode.Id)

			var order models.OrderNode
			o = orm.NewOrm()
			qs = o.QueryTable("order_node")
			err = qs.Filter("id", nodes.OrderNode.Id).One(&order)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(order)
	*/

	// 创建 user-service 微服务的客户端
	client := vps.NewVpsService(serviceName, srv.Client())

	//name := "Ethan"

	retrys := 0
	for {
		retrys++
		fmt.Println("retrys:", retrys)
		if retrys > 10 {
			log.Printf("call Create error: %v", "timeout")
		}
		resp, err := client.Call(context.Background(), &vps.Request{
			UserId: 1,
			Id:     18,
		})
		if err == nil {
			log.Println("created: ", "success call")
			fmt.Println(resp)
			break
		}
		time.Sleep(time.Second)
	}

	/*
		resp, err := client.GetNode(context.Background(), &vps.Request{
			UserId: 1,
			Id:     111,
		})
		if err != nil {
			log.Printf("call Create error: %v", err)
		} else {
			log.Println("created: ", "success call")
		}
		fmt.Println(resp)
	*/

	/*
		var id int64
		id = 1
		resp, err := client.NewNode(context.Background(), &vps.Request{
			UserId: 1,
			Id:     id,
		})
		if err != nil {
			log.Printf("call Create error: %v", err)
		} else {
			log.Println("created: ", "success call")
		}
		fmt.Println(resp)
	*/
}

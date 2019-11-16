package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"log"

	"strconv"

	"github.com/micro/go-micro/broker"
	"gopkg.in/mgo.v2"
	"vircle.io/mnhost/config"
	logPB "vircle.io/mnhost/interface/out/log"
	nodePB "vircle.io/mnhost/interface/out/nodemsg"
	"vircle.io/mnhost/model"
	"vircle.io/mnhost/utils"
)

// 微服务服务端 struct handler 必须实现 protobuf 中定义的 rpc 方法
// 实现方法的传参等可参考生成的 consignment.pb.go
type Nodemsg struct {
	session *mgo.Session
	Broker  broker.Broker
}

const service = "nodemsg"

var (
	topic       string
	serviceName string
	version     string
)

func init() {
	topic = config.GetBrokerTopic("log")
	serviceName = config.GetServiceName(service)
	version = config.GetVersion(service)
	if version == "" {
		version = "latest"
	}
}

func GetHandler(session *mgo.Session, bk broker.Broker) *Vps {
	return &Vps{
		session: session,
		Broker:  bk,
	}
}

// 发送vps
func (e *Nodemsg) pubMsg(userID, topic string, msgId int64) error {
	fmt.Println("start msg publish:", userID, "--", topic)
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

	if err := e.Broker.Publish(topic, data); err != nil {
		log.Printf("[pub] failed: %v\n", err)
	}
	fmt.Println("success msg publish:", userID, "--", topic)
	return nil
}



func (e *Nodemsg) pubLog(userID, method, msg string) error {
	fmt.Println("start log publish")
	logPB := logPB.Log{
		Method: method,
		Origin: serviceName,
		Msg:    msg,
	}
	body, err := json.Marshal(logPB)
	if err != nil {
		return err
	}

	data := &broker.Message{
		Header: map[string]string{
			"user_id": userID,
		},
		Body: body,
	}

	if err := e.Broker.Publish(topic, data); err != nil {
		log.Printf("[pub] failed: %v\n", err)
	}
	fmt.Println("success log publish")
	return nil
}

// Call is a single request handler called via client.Call or the generated client code
func (e *Nodemsg) Call(ctx context.Context, req *vps.Request, rsp *vps.Response) error {
	log.Info("Received Vps.Call request")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

		return nil

	}
	e.pubLog("123456", "neworder", "abcd1234567890cccc")
	e.pubMsg("123456", "abcd1234567890cccc", "Vircle.Topic.newOrder")

	return nil
}


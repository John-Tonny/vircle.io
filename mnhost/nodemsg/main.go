package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	vps "vircle.io/mnhost/interface/proto/vps"

	json "github.com/json-iterator/go"
	"github.com/micro/go-micro/broker"

	//"github.com/astaxie/beego/orm"

	"vircle.io/mnhost/common"
	"vircle.io/mnhost/config"

	//"vircle.io/mnhost/model"
	"vircle.io/mnhost/utils"

	logPB "vircle.io/mnhost/interface/out/log"
	nodePB "vircle.io/mnhost/interface/out/nodemsg"
)

const cservice = "nodemsg"
const oservice = "vps"
const topic_newnode_success = "Vircle.Mnhost.Topic.NodeNew.Success"
const topic_newnode_fail = "Vircle.Mnhost.Topic.NodeNew.Fail"
const topic_newnode_start = "Vircle.Mnhost.Topic.NodeNew.Start"
const topic_newnode_stop = "Vircle.Mnhost.Topic.NodeNew.Stop"

const topic_delnode_success = "Vircle.Mnhost.Topic.NodeDel.Success"
const topic_delnode_fail = "Vircle.Mnhost.Topic.NodeDel.Fail"
const topic_delnode_start = "Vircle.Mnhost.Topic.NodeDel.Start"
const topic_delnode_stop = "Vircle.Mnhost.Topic.NodeDel.Stop"

var (
	topic        string
	cserviceName string
	oserviceName string
	gBroker      broker.Broker
)

func init() {
	cserviceName = config.GetServiceName(cservice)
	oserviceName = config.GetServiceName(oservice)
	topic = config.GetBrokerTopic("log")
}

func main() {
	log.Println("nodemsg service start")

	srv := common.GetMicroServer(cservice)

	bk := srv.Server().Options().Broker
	gBroker = bk

	// 将实现服务端的 API 注册到服务端
	//vps.RegisterVpsHandler(srv.Server(), handler.GetHandler(session, bk))

	// 这里订阅了 一个 topic, 并提供接口处理
	_, err := bk.Subscribe(topic_newnode_start, nodeNewStart)
	if err != nil {
		log.Fatalf("new node start error: %v\n", err)
	}

	_, err = bk.Subscribe(topic_newnode_success, nodeNewSuccess)
	if err != nil {
		log.Fatalf("new node success error: %v\n", err)
	}

	_, err = bk.Subscribe(topic_newnode_fail, nodeNewFail)
	if err != nil {
		log.Fatalf("new node fail error: %v\n", err)
	}

	_, err = bk.Subscribe(topic_newnode_stop, nodeNewStop)
	if err != nil {
		log.Fatalf("new node stop error: %v\n", err)
	}

	_, err = bk.Subscribe(topic_delnode_start, nodeDelStart)
	if err != nil {
		log.Fatalf("del node start error: %v\n", err)
	}

	_, err = bk.Subscribe(topic_delnode_success, nodeDelSuccess)
	if err != nil {
		log.Fatalf("del node success error: %v\n", err)
	}

	_, err = bk.Subscribe(topic_delnode_fail, nodeDelFail)
	if err != nil {
		log.Fatalf("del node fail error: %v\n", err)
	}

	_, err = bk.Subscribe(topic_delnode_stop, nodeDelStop)
	if err != nil {
		log.Fatalf("del node stop error: %v\n", err)
	}

	if err = srv.Run(); err != nil {
		log.Fatalf("srv run error: %v\n", err)
	}
}

func nodeNewStart(pub broker.Event) error {
	log.Println("new node start")
	userId := pub.Message().Header["user_id"]
	method := "nodeNewStart"
	var msg *nodePB.NodeMsg
	if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
		fmt.Println("json:", err)
		pubErrMsg(userId, method, utils.JSON_DATAERR, "", topic_newnode_fail)
		return err
	}
	id := (*msg).Id
	log.Println("msg:", userId, id)

	srv := common.GetMicroClient(oservice)

	fmt.Println(oservice, ":", oserviceName)

	// 创建 user-service 微服务的客户端
	client := vps.NewVpsService(oserviceName, srv.Client())

	iuserId, err := strconv.ParseInt(userId, 10, 64)
	sId := strconv.FormatInt(id, 10)
	log.Println("type:", reflect.TypeOf(iuserId), reflect.TypeOf(sId))
	log.Println(userId, id)
	if err != nil {
		log.Println("aaa")
		pubErrMsg(userId, method, utils.RECODE_DATAERR, sId, topic_newnode_fail)
		return err
	}

	retrys := 0
	for {
		retrys++
		fmt.Println("retrys:", retrys)
		if retrys > 10 {
			pubErrMsg(userId, method, utils.TIMEOUT_VPS, sId, topic_newnode_fail)
			return errors.New("timeout")
		}
		resp, err := client.NewNode(context.Background(), &vps.Request{
			UserId: iuserId,
			Id:     id,
		})
		if err == nil {
			fmt.Println("new node:", resp)
			if resp.Errno == utils.RECODE_OK {
				pubMsg(userId, method, topic_newnode_stop, id)
				return nil
			}
		}
		fmt.Println(err)
		time.Sleep(time.Second)
	}

	pubErrMsg(userId, method, utils.RECODE_SERVERERR, sId, topic_newnode_fail)
	/*
		var order models.OrderNode
		o := orm.NewOrm()
		qs := o.QueryTable("order_node")
		err := qs.Filter("id", orderId).One(&order)
		if err != nil {
			pubErrMsg(userID, method, utils.RECODE_NODATA, "", "Vircle.Mnhost.Topic.ErrOrder")
			return err
		}

		fmt.Println("new order:", order.Id)
	*/

	log.Println("new node start finish")
	return nil
}

func nodeNewSuccess(pub broker.Event) error {
	log.Println("new node success")
	var msg *nodePB.NodeMsg
	if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
		return err
	}
	userId := pub.Message().Header["user_id"]
	orderId := (*msg).Id
	log.Println("msg:", userId, orderId)
	log.Println("new node success finish")
	return nil
}

func nodeNewFail(pub broker.Event) error {
	log.Printf("new node fail")
	var msg *nodePB.NodeErrMsg
	if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
		return err
	}
	userId := pub.Message().Header["user_id"]
	log.Println("msg:", userId, msg)
	log.Printf("new node fail finish")
	return nil
}

func nodeNewStop(pub broker.Event) error {
	log.Printf("new node stop")
	var msg *nodePB.NodeMsg
	if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
		return err
	}

	userId := pub.Message().Header["user_id"]
	orderId := (*msg).Id
	log.Println("msg:", userId, orderId)
	log.Printf("new node stop finish")
	return nil
}

func nodeDelStart(pub broker.Event) error {
	log.Println("del node start")
	userId := pub.Message().Header["user_id"]
	method := "nodeDelStart"
	var msg *nodePB.NodeMsg
	if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
		pubErrMsg(userId, method, utils.JSON_DATAERR, "", topic_newnode_fail)
		return err
	}
	id := (*msg).Id
	log.Println("msg:", userId, id)

	srv := common.GetMicroClient(oservice)

	fmt.Println(oservice, ":", oserviceName)

	// 创建 user-service 微服务的客户端
	client := vps.NewVpsService(oserviceName, srv.Client())

	iuserId, err := strconv.ParseInt(userId, 10, 64)
	sId := strconv.FormatInt(id, 10)
	if err != nil {
		pubErrMsg(userId, method, utils.RECODE_DATAERR, sId, topic_newnode_fail)
		return err
	}
	resp, err := client.DelNode(context.Background(), &vps.Request{
		UserId: iuserId,
		Id:     id,
	})
	if err != nil {
		pubErrMsg(userId, method, utils.RECODE_SERVERERR, sId, topic_newnode_fail)
		return err
	}
	fmt.Println("del node:", resp)

	if resp.Errno == utils.RECODE_OK {
		pubMsg(userId, method, topic_delnode_stop, id)
		return nil
	}
	pubErrMsg(userId, method, utils.RECODE_SERVERERR, sId, topic_newnode_fail)
	/*
		var order models.OrderNode
		o := orm.NewOrm()
		qs := o.QueryTable("order_node")
		err := qs.Filter("id", orderId).One(&order)
		if err != nil {
			pubErrMsg(userID, method, utils.RECODE_NODATA, "", "Vircle.Mnhost.Topic.ErrOrder")
			return err
		}

		fmt.Println("new order:", order.Id)
	*/

	log.Println("del node start finish")
	return nil
}

func nodeDelSuccess(pub broker.Event) error {
	log.Println("del node success")
	var msg *nodePB.NodeMsg
	if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
		return err
	}
	userId := pub.Message().Header["user_id"]
	nodeId := (*msg).Id
	log.Println("msg:", userId, nodeId)
	log.Println("del node success finish")
	return nil
}

func nodeDelFail(pub broker.Event) error {
	log.Printf("del node fail")
	var msg *nodePB.NodeErrMsg
	if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
		return err
	}
	userId := pub.Message().Header["user_id"]
	log.Println("msg:", userId, msg)
	log.Printf("del node fail finish")
	return nil
}

func nodeDelStop(pub broker.Event) error {
	log.Printf("del node stop")
	var msg *nodePB.NodeMsg
	if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
		return err
	}

	userId := pub.Message().Header["user_id"]
	nodeId := (*msg).Id
	log.Println("msg:", userId, nodeId)
	log.Printf("del node stop finish")
	return nil
}

// 发送
func pubMsg(userID, method, topic string, msgId int64) error {
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

	if err := gBroker.Publish(topic, data); err != nil {
		log.Printf("[pub] failed: %v\n", err)
	}
	fmt.Println("success msg publish:", userID, "--", topic)
	return nil
}

func pubErrMsg(userID, method, errno, msg, topic string) error {
	fmt.Println("start err msg publish:", userID, "--", topic)
	errmsg := nodePB.NodeErrMsg{
		Method: method,
		Origin: oserviceName,
		Errno:  errno,
		Errmsg: utils.RecodeText(errno),
		Msg:    msg,
	}
	body, err := json.Marshal(errmsg)
	if err != nil {
		return err
	}

	data := &broker.Message{
		Header: map[string]string{
			"user_id": userID,
		},
		Body: body,
	}

	if err := gBroker.Publish(topic, data); err != nil {
		log.Printf("[pub] failed: %v\n", err)
	}
	fmt.Println("success err msg publish:", userID, "--", topic)
	return nil
}

func pubLog(userID, method, msg string) error {
	fmt.Println("start log publish")
	logPB := logPB.Log{
		Method: method,
		Origin: oserviceName,
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

	if err := gBroker.Publish(topic, data); err != nil {
		log.Printf("[pub] failed: %v\n", err)
	}
	fmt.Println("success log publish")
	return nil
}

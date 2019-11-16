package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"reflect"
	"time"

	"github.com/go-ini/ini"
	"github.com/pytool/ssh"

	"github.com/dynport/gossh"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	// "github.com/micro/go-micro/util/log"

	//"github.com/micro/go-micro"

	vps "vircle.io/mnhost/interface/proto/vps"

	uec2 "github.com/John-Tonny/micro/vps/amazon"

	log "github.com/sirupsen/logrus"

	"strconv"

	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"

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
type Vps struct {
	session *mgo.Session
	Broker  broker.Broker
}

type VpsInfo struct {
	instanceId      string
	regionName      string
	allocationId    string
	allocationState bool
	publicIp        string
	volumeId        string
	volumeState     bool
}

type Volume struct {
	Mountpoint string `json:"Mountpoint"`
	Name       string `json:"Name"`
}

func logInit() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	file := time.Now().Format("20060102") + ".txt" //文件名
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		panic(err)
	}
	log.SetOutput(logFile)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

const service = "vps"
const topic_newnode_success = "Vircle.Mnhost.Topic.NodeNew.Success"
const topic_newnode_fail = "Vircle.Mnhost.Topic.NodeNew.Fail"
const topic_newnode_start = "Vircle.Mnhost.Topic.NodeNew.Start"

const topic_delnode_success = "Vircle.Mnhost.Topic.NodeDel.Success"
const topic_delnode_fail = "Vircle.Mnhost.Topic.NodeDel.Fail"
const topic_delnode_start = "Vircle.Mnhost.Topic.NodeDel.Start"
const ssh_password = "vpub$999000"
const rpc_user = "vpub"
const rpc_password = "vpub999000"
const port_from = 19900
const port_to = 20000

var (
	topic       string
	serviceName string
	version     string
)

func init() {
	logInit()
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
func (e *Vps) pubMsg(userID, topic string, msgId int64) error {
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

func (e *Vps) pubErrMsg(userID, method, errno, msg, topic string) error {
	fmt.Println("start err msg publish:", userID, "--", topic)
	errmsg := nodePB.NodeErrMsg{
		Method: method,
		Origin: serviceName,
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

	if err := e.Broker.Publish(topic, data); err != nil {
		log.Printf("[pub] failed: %v\n", err)
	}
	fmt.Println("success err msg publish:", userID, "--", topic)
	return nil
}

func (e *Vps) pubLog(userID, method, msg string) error {
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
func (e *Vps) Call(ctx context.Context, req *vps.Request, rsp *vps.Response) error {
	log.Info("Received Vps.Call request")
	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	if req.UserId == 1 {
		var node models.Node
		o := orm.NewOrm()
		qs := o.QueryTable("node")
		fmt.Println("nodeId:", req.Id)
		err := qs.Filter("id", req.Id).One(&node)
		if err != nil {
			rsp.Errno = utils.RECODE_NODATA
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}

		jnode, err := json.Marshal(node)
		fmt.Println(jnode)
		rsp.Mix = jnode

		e.pubLog("1", "new node", topic_delnode_start)
		e.pubMsg("1", topic_delnode_start, req.Id)
	} else {
		var order models.OrderNode
		o := orm.NewOrm()
		qs := o.QueryTable("order_node")
		fmt.Println("orderId:", req.Id)
		err := qs.Filter("id", req.Id).One(&order)
		if err != nil {
			rsp.Errno = utils.RECODE_NODATA
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}

		jorder, err := json.Marshal(order)
		fmt.Println(jorder)
		rsp.Mix = jorder

		e.pubLog("1", "new node", topic_newnode_start)
		e.pubMsg("1", topic_newnode_start, req.Id)
	}
	return nil
}

func (e *Vps) NewNode(ctx context.Context, req *vps.Request, rsp *vps.Response) error {
	log.Info("Received Vps.NewNode request")
	
	orderId := req.Id
	userId := strconv.Itoa(int(req.UserId))
	fmt.Println("orderId:", orderId, reflect.TypeOf(orderId))

	var order models.OrderNode
	o := orm.NewOrm()
	qs := o.QueryTable("order_node")
	err := qs.Filter("id", req.Id).One(&order)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	/*if order.Status == models.ORDER_STATUS_PROCESSING {
		rsp.Errno = utils.ORDER_PROCESSING
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}*/

	/*o = orm.NewOrm()
	order.Status = models.ORDER_STATUS_PROCESSING
	_, err = o.Update(&order)
	if err != nil {
		fmt.Println("order_node update失败", err)
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}*/

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	jorder, err := json.Marshal(order)
	rsp.Mix = jorder

	go e.processNewNode(userId, int64(order.Id))

	return nil
	/*
		log.Info("Received Vps.NewNode request")

		orderId := req.Id
		fmt.Println("orderId:", orderId)

		var order models.OrderNode
		o := orm.NewOrm()
		qs := o.QueryTable("order_node")
		err := qs.Filter("id", orderId).One(&order)
		if err != nil {
			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}

		if order.status == models.ORDER_STATUS_PROCESSING {
			rsp.Errno = utils.ORDER_PROCESSING
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
		var order1 models.OrderNode
		order1.status = models.ORDER_STATUS_PROCESSING

		o = orm.NewOrm()
		vps.UsableNodes = vps.UsableNodes - 1
		_, err = o.Update(&vps)
		if err != nil {
			fmt.Println("vps update失败", err)
			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
		fmt.Println("update vps finish")

		var nvpsInfo VpsInfo
		var vpsInfo *VpsInfo
		var vps models.Vps
		nvps := models.Vps{}
		o = orm.NewOrm()
		qs = o.QueryTable("vps")
		err = qs.Filter("usable_nodes__gt", 0).One(&vps)
		retrys := 0
		if err != nil {
			fmt.Println("vps node 查询失败", err)
			for { //循环
				retrys++
				if retrys >= 20 {
					fmt.Println("retrys:", retrys)
					rsp.Errno = "new vps timeout"
					rsp.Errmsg = utils.RecodeText(rsp.Errno)
					return nil
				}
				vpsInfo, err = NewVps("", "", "", 1, nvpsInfo)
				fmt.Println("retrys:", retrys)
				if err != nil {
					time.Sleep(5 * time.Second)
					continue
				}
				nvpsInfo.allocationId = vpsInfo.allocationId
				nvpsInfo.allocationState = vpsInfo.allocationState
				nvpsInfo.instanceId = vpsInfo.instanceId
				nvpsInfo.publicIp = vpsInfo.publicIp
				nvpsInfo.regionName = vpsInfo.regionName
				nvpsInfo.volumeId = vpsInfo.volumeId
				nvpsInfo.volumeState = vpsInfo.volumeState
				break
			}
			fmt.Println("success new vps:", vpsInfo.instanceId)

			log.Info(vpsInfo)
			nvps.AllocateId = vpsInfo.allocationId
			nvps.InstanceId = vpsInfo.instanceId
			nvps.VolumeId = vpsInfo.volumeId
			nvps.ProviderName = "amazon"
			nvps.Cores = 1
			nvps.Memory = 1
			nvps.KeyPairName = "vcl-keypair"
			nvps.MaxNodes = 5
			nvps.UsableNodes = 5
			nvps.SecurityGroupName = "vcl-mngroup"
			nvps.RegionName = vpsInfo.regionName
			nvps.IpAddress = vpsInfo.publicIp
			o = orm.NewOrm()
			_, err = o.Insert(&nvps)
			if err != nil {
				fmt.Println("vps插入失败", err)
				rsp.Errno = utils.RECODE_DBERR
				rsp.Errmsg = utils.RecodeText(rsp.Errno)
				return nil
			}
			o = orm.NewOrm()
			qs = o.QueryTable("vps")
			err = qs.Filter("usable_nodes__gt", 0).One(&vps)
			if err != nil {
				rsp.Errno = utils.RECODE_DBERR
				rsp.Errmsg = utils.RecodeText(rsp.Errno)
				return err
			}
		}
		fmt.Printf("%d-%s-%s\n", vps.UsableNodes, order.Coin, vps.IpAddress)

		retrys = 0
		for { //循环
			retrys++
			if retrys >= 5 {
				fmt.Println("retrys:", retrys)
				rsp.Errno = "new node timeout"
				rsp.Errmsg = utils.RecodeText(rsp.Errno)
				return nil
			}
			err = NewNode(vps.UsableNodes, order.Coin, vps.IpAddress, ssh_password)
			fmt.Println("retrys:", retrys)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			break
		}
		fmt.Println("success new node")

		node := models.Node{}
		node.Coin = order.Coin
		node.User = order.User
		node.Vps = &vps
		node.OrderNode = &order
		node.Port = 20000 - vps.UsableNodes*2
		o = orm.NewOrm()
		nodeId, err := o.Insert(&node)
		if err != nil {
			fmt.Println("vps插入失败", err)
			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
		fmt.Println("success insert node")

		o = orm.NewOrm()
		vps.UsableNodes = vps.UsableNodes - 1
		_, err = o.Update(&vps)
		if err != nil {
			fmt.Println("vps update失败", err)
			rsp.Errno = utils.RECODE_DBERR
			rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
		fmt.Println("update vps finish")

		rsp.Errno = utils.RECODE_OK
		rsp.Errmsg = utils.RecodeText(rsp.Errno)

		userId := strconv.Itoa(order.User.Id)
		fmt.Println("userId:", userId)
		//e.pubLog(userId, "NewNode", vpsInfo.instanceId)
		//e.pubMsg(userId, "Vircle.Mnhost.Topic.FinishNode", nodeId)

		node1, err := json.Marshal(node)
		fmt.Println(node1)
		rsp.Mix = node1
		fmt.Println("success order:", nodeId)
		return nil
	*/
}

func (e *Vps) GetNode(ctx context.Context, req *vps.Request, rsp *vps.Response) error {
	log.Info("Received get node request")

	nodeId := req.Id
	userId := strconv.Itoa(int(req.UserId))
	fmt.Println("params:", userId, nodeId)

	var node models.Node
	o := orm.NewOrm()
	qs := o.QueryTable("node")
	err := qs.Filter("id", nodeId).One(&node)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	jnode, err := json.Marshal(node)
	rsp.Mix = jnode

	return nil
}

func (e *Vps) DelNode(ctx context.Context, req *vps.Request, rsp *vps.Response) error {
	log.Info("Received del node request")

	nodeId := req.Id
	userId := strconv.Itoa(int(req.UserId))
	fmt.Println("params:", userId, nodeId)

	var node models.Node
	o := orm.NewOrm()
	qs := o.QueryTable("node")
	err := qs.Filter("id", nodeId).One(&node)
	if err != nil {
		rsp.Errno = utils.RECODE_DBERR
		rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	rsp.Errno = utils.RECODE_OK
	rsp.Errmsg = utils.RecodeText(rsp.Errno)

	jnode, err := json.Marshal(node)
	rsp.Mix = jnode

	go e.processDelNode(userId, int64(node.Id))

	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Vps) Stream(ctx context.Context, req *vps.StreamingRequest, stream vps.Vps_StreamStream) error {
	log.Info("Received Vps.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Info("Responding: %d", i)
		if err := stream.Send(&vps.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Vps) PingPong(ctx context.Context, stream vps.Vps_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Info("Got ping %v", req.Stroke)
		if err := stream.Send(&vps.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

func NewVps(imageId, zone, instanceType string, volumeSize int64, nvpsInfo VpsInfo) (*VpsInfo, error) {
	groupName := "vcl-mngroup"
	groupDesc := "basic masternode group"
	keyPairName := "vcl-keypair"
	deviceName := "xvdk"

	/*nvpsInfo.instanceId = "i-001da7ecc8a29947b"
	nvpsInfo.allocationState = true
	nvpsInfo.allocationId = "eipalloc-0af3980b97808b4a9"
	nvpsInfo.volumeId = "vol-001f373aea56506d6"*/

	if len(instanceType) == 0 {
		instanceType = "t2.micro"
	}

	if volumeSize == 0 {
		volumeSize = 20
	}

	var vpsInfo VpsInfo
	vpsInfo.allocationState = false
	vpsInfo.volumeState = false

	c, err := uec2.NewEc2Client(zone, "test-account")
	if err != nil {
		return &vpsInfo, err
	}
	log.Info("client")

	var securityGroupId string
	groupResult, err := c.GetDescribeSecurityGroupsFromName([]string{groupName})
	if err != nil {
		ipPermissions := GetIpPermission()
		securityGroupId, err = c.CreateSecurityGroups(groupName, groupDesc, ipPermissions)
		if err != nil {
			return &vpsInfo, err
		}
		log.Info(securityGroupId)
	} else {
		securityGroupId = aws.StringValue(groupResult.SecurityGroups[0].GroupId)
		fmt.Println(securityGroupId)
	}

	log.Info("group")

	_, err = c.GetDescribeKeyPairs([]string{keyPairName})
	if err != nil {
		_, err := c.CreateKeyPairs(keyPairName)
		if err != nil {
			return &vpsInfo, err
		}
	}
	log.Info("keypair")

	log.Info("instance")
	retrys := 0
	regionName := ""
	instanceId := nvpsInfo.instanceId
	fmt.Println("instanceId:", instanceId)
	for { //循环
		retrys++
		if retrys >= 30 {
			fmt.Println("retrys:", retrys)
			err = errors.New("instance timeout")
			return &vpsInfo, err
		}
		result, err := c.GetDescribeInstance([]string{instanceId})
		fmt.Println("get instance retrys:", retrys)
		if err != nil {
			instanceId, err = c.CreateInstances(imageId, instanceType, keyPairName, securityGroupId)
			fmt.Println("error:", err)
			continue
		}
		regionName = aws.StringValue(result.Reservations[0].Instances[0].Placement.AvailabilityZone)
		state := aws.StringValue(result.Reservations[0].Instances[0].State.Name)
		if state == "running" {
			fmt.Println("vps state:", state)
			break
		} else {
			c.StartInstance(instanceId)
			time.Sleep(time.Second)
			continue
		}
		time.Sleep(time.Second)
	}
	vpsInfo.instanceId = instanceId
	vpsInfo.regionName = regionName
	fmt.Println("create instance")

	allocationId := ""
	publicIp := ""
	if len(nvpsInfo.allocationId) == 0 {
		publicIp, allocationId, err = c.AllocateAddresss(instanceId)
		if err != nil {
			return &vpsInfo, err
		}
	}
	log.Info(publicIp)
	vpsInfo.allocationId = allocationId
	vpsInfo.publicIp = publicIp

	if nvpsInfo.allocationState == false {
		_, err = c.AssociateAddresss(instanceId, allocationId)
		if err != nil {
			return &vpsInfo, err
		}
	}
	log.Info("associate")
	vpsInfo.allocationState = true

	volumeId := ""
	if len(nvpsInfo.volumeId) == 0 {
		fmt.Println("create volume")
		volumeId, err = c.CreateVolumes(regionName, volumeSize)
		if err != nil {
			return &vpsInfo, err
		}
	} else {
		volumeId = nvpsInfo.volumeId
	}
	vpsInfo.volumeId = volumeId

	if nvpsInfo.volumeState == false {
		retrys = 0
		for { //循环
			retrys++
			if retrys >= 10 {
				err = errors.New("volume timeout")
				return &vpsInfo, err
			}
			vResult, err := c.GetDescribeVolumes([]string{volumeId})
			fmt.Println("volume retrys:", retrys)
			if err == nil {
				vState := aws.StringValue(vResult.Volumes[0].State)
				if vState == "available" {
					fmt.Println("volume state:", vState)
					break
				}
			}
			time.Sleep(time.Second)
		}
		fmt.Println("create volumn avalid", instanceId, volumeId, deviceName)
		_, err = c.AttachVolumes(instanceId, volumeId, deviceName)
		if err != nil {
			fmt.Println("ass volume err:", err)
			return &vpsInfo, err
		}
	}
	fmt.Println("success update vps")
	vpsInfo.volumeState = true

	return &vpsInfo, nil
}

func MakeLogger(prefix string) gossh.Writer {
	return func(args ...interface{}) {
		log.Println((append([]interface{}{prefix}, args...))...)
	}
}

func NewNode(vpsId int, coinName string, ipAddress string, password string) error {
	var coin models.Coin
	o := orm.NewOrm()
	qs := o.QueryTable("coin")
	err := qs.Filter("name", coinName).One(&coin)
	if err != nil {
		fmt.Println("coin node 查询失败", err)
	}
	log.Info("new node")

	client := gossh.New(ipAddress, "root")

	// my default agent authentication is used. use
	client.SetPassword(ssh_password)
	// for password authentication
	//client.DebugWriter = MakeLogger("DEBUG")
	client.InfoWriter = MakeLogger("INFO ")
	client.ErrorWriter = MakeLogger("ERROR")

	defer client.Close()

	log.Info("get volumn")
	fmt.Println("vpsId:", vpsId)
	rpcPort, port, err := getRpcPort(vpsId)
	if err != nil {
		fmt.Println("port error:", err)
		return err
	}
	fmt.Println(rpcPort, port)

	volumeName := "mn" + strconv.Itoa(rpcPort)
	fmt.Println(volumeName)
	var rsp *gossh.Result
	cmd := "docker volume inspect " + volumeName
	rsp, err = sshCmd(client, cmd)
	if err != nil {
		log.Info("create volume")
		cmd = "docker volume create --name=" + volumeName
		rsp, err = sshCmd(client, cmd)
		if err != nil {
			fmt.Println("bbb")
			client.ErrorWriter(err.Error())
			client.ErrorWriter("STDOUT: " + rsp.Stdout())
			client.ErrorWriter("STDERR: " + rsp.Stderr())
			return err
		}
		time.Sleep(5 * time.Second)
		cmd = "docker volume inspect " + volumeName
		rsp, err = sshCmd(client, cmd)
		if err != nil {
			return err
		}
		fmt.Println("search volume:", rsp.Stdout())
	}
	fmt.Println(rsp.Stdout())
	var part []Volume
	if err = json.Unmarshal([]byte(rsp.Stdout()), &part); err != nil {
		return err
	}
	mountPoint := part[0].Mountpoint
	fmt.Println(mountPoint)

	coinPath := mountPoint + "/" + coin.Path
	fmt.Println(coinPath)
	cmd = "cd " + coinPath
	rsp, err = sshCmd(client, cmd)
	if err != nil {
		cmd = "mkdir " + coinPath
		_, err = sshCmd(client, cmd)
		if err != nil {
			return err
		}
		log.Info("mkdir")
	}
	log.Info("cd dir")

	localFile := "/tmp/" + coin.Conf
	fmt.Println(localFile)
	err = WriteConf(localFile, port)
	if err != nil {
		return err
	}
	log.Info("write conf", localFile)

	err = UploadFile(ipAddress, "root", password, localFile, "/tmp/")
	if err != nil {
		return err
	}
	log.Info("upload conf")

	cmd = "mv /tmp/" + coin.Conf + " " + coinPath
	_, err = sshCmd(client, cmd)
	if err != nil {
		return err
	}
	fmt.Println("mv conf")

	cmd = "docker images | grep " + coinName
	rsp, err = sshCmd(client, cmd)
	if err != nil {
		localFile = "/root/vpub-vircle-0.1.tar"
		err = UploadFile(ipAddress, "root", password, localFile, "/tmp/")
		if err != nil {
			return err
		}

		cmd = "docker load  --input /tmp/vpub-vircle-0.1.tar"
		rsp, err = sshCmd(client, cmd)
		if err != nil {
			client.ErrorWriter(err.Error())
			client.ErrorWriter("STDOUT: " + rsp.Stdout())
			client.ErrorWriter("STDERR: " + rsp.Stderr())
			return err
		}
		fmt.Println("load docker")
	}
	fmt.Println("cp docker")

	srpcPort := strconv.Itoa(rpcPort)
	sport := strconv.Itoa(port)
	cmd = "docker run -v " + volumeName + ":/" + coinName + " --name=" + volumeName + " -d -p " + sport + ":" + sport + " -p " + srpcPort + ":" + srpcPort + " " + coin.Docker
	fmt.Println(cmd)
	rsp, err = sshCmd(client, cmd)
	if err != nil {
		return err
	}
	log.Info("start docker")

	return nil
}

func GetIpPermission() []*ec2.IpPermission {
	ipPermissions := []*ec2.IpPermission{
		(&ec2.IpPermission{}).
			SetIpProtocol("tcp").
			SetFromPort(80).
			SetToPort(80).
			SetIpRanges([]*ec2.IpRange{
				{CidrIp: aws.String("0.0.0.0/0")},
			}),
		(&ec2.IpPermission{}).
			SetIpProtocol("tcp").
			SetFromPort(22).
			SetToPort(22).
			SetIpRanges([]*ec2.IpRange{
				(&ec2.IpRange{}).
					SetCidrIp("0.0.0.0/0"),
			}),
		(&ec2.IpPermission{}).
			SetIpProtocol("tcp").
			SetFromPort(port_from).
			SetToPort(port_to).
			SetIpRanges([]*ec2.IpRange{
				(&ec2.IpRange{}).
					SetCidrIp("0.0.0.0/0"),
			}),
	}
	return ipPermissions
}

func SecurityGroupIsExist(groupName string, groupResult *ec2.DescribeSecurityGroupsOutput) string {
	for _, group := range groupResult.SecurityGroups {
		if aws.StringValue(group.GroupName) == groupName {
			return aws.StringValue(group.GroupId)
		}
	}
	return ""
}

func Testuser(name string, password string, mobile string) error {
	fmt.Println(" 注册服务  PostRet  /api/v1.0/users")
	//1.初始化返回值
	//rsp.Errno = utils.RECODE_OK
	//rsp.Errmsg = utils.RecodeText(rsp.Errno)
	/****2.连接redis**/
	bm, err := utils.RedisOpen(utils.G_server_name, utils.G_redis_addr, utils.G_redis_port, utils.G_redis_dbnum)
	if err != nil {
		log.Println("redis连接失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	/****3.从redis中获取短信验证码**/
	value := bm.Get(mobile)
	value_string, _ := redis.String(value, nil)
	/****4.检查短信验证码是否正确**/
	if value_string == "aaaa" {
		fmt.Println("短信验证码错误", value_string, "aaaa")
		//rsp.Errno = utils.RECODE_DATAERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	/****5.对接收到的密码进行加密**/
	user := models.User{}
	user.Password_hash = utils.Getmd5string(password)
	user.Mobile = mobile
	user.Name = name
	/****6.插入数据到数据库中**/
	o := orm.NewOrm()
	id, err := o.Insert(&user)
	if err != nil {
		fmt.Println("用户数据注册插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	/****7.生成sessionid**/
	sessionid := utils.Getmd5string(mobile + password + strconv.Itoa(int(time.Now().UnixNano())))
	/****8.通过sessionid将数据返回redis**/
	//rsp.Sessionid = sessionid
	bm.Put(sessionid+"user_id", id, time.Second*600)
	bm.Put(sessionid+"mobile", user.Mobile, time.Second*600)
	bm.Put(sessionid+"name", user.Name, time.Second*600)
	return nil
}

func Testorder(mobile string) error {
	var user models.User
	o := orm.NewOrm()
	qs := o.QueryTable("user")
	err := qs.Filter("mobile", mobile).One(&user)
	if err != nil {
		fmt.Println("用户名查询失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	fmt.Println(" 注册服务  PostRet  /api/v1.0/users")
	//1.初始化返回值
	//rsp.Errno = utils.RECODE_OK
	//rsp.Errmsg = utils.RecodeText(rsp.Errno)
	/****2.连接redis**/
	order := models.OrderNode{}
	order.Coin = "vircle"
	order.Alias = "vircle_mn03"
	order.Txid = "2343243242-23432432"
	order.OutputIndex = 0
	order.RewardAddress = "abcdef123456"

	order.Begin_date, _ = time.ParseInLocation("2006-01-02 15:04:05", "2019-11-05 16:00:00", time.Local)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(order.Begin_date)

	order.End_date, _ = time.ParseInLocation("2006-01-02 15:04:05", "2019-12-05 15:59:59", time.Local)
	order.Period = "1month"
	order.Amount = 30
	order.Status = models.ORDER_STATUS_WAIT_PAYMENT

	order.User = &user

	/****6.插入数据到数据库中**/
	o = orm.NewOrm()
	id, err := o.Insert(&order)
	if err != nil {
		fmt.Println("数据插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	fmt.Println("order success", id, err)
	return nil
}

func Testnode(orderid string) error {
	var order models.OrderNode
	o := orm.NewOrm()
	qs := o.QueryTable("order_node")
	err := qs.Filter("id", orderid).One(&order)
	if err != nil {
		fmt.Println("order node 查询失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	fmt.Println(" node服务   /api/v1.0/orders")
	fmt.Println(order.Id)

	var vpsInfo VpsInfo
	var vps models.Vps
	o = orm.NewOrm()
	qs = o.QueryTable("vps")
	err = qs.Filter("usable_nodes__gt", 0).One(&vps)
	log.Info("aaa")
	if err != nil {
		fmt.Println("vps node 查询失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		vpsInfo, err := NewVps("", "", "", 1, vpsInfo)
		if err != nil {
			return err
		}
		log.Info(vpsInfo)
		vps1 := models.Vps{}
		vps1.AllocateId = vpsInfo.allocationId
		vps1.InstanceId = vpsInfo.instanceId
		vps1.VolumeId = vpsInfo.volumeId
		vps1.ProviderName = "amazon"
		vps1.Cores = 1
		vps1.Memory = 1
		vps1.KeyPairName = "vcl-keypair"
		vps1.MaxNodes = 10
		vps1.UsableNodes = 10
		vps1.SecurityGroupName = "vcl-mngroup"
		vps1.RegionName = vpsInfo.regionName
		vps1.IpAddress = vpsInfo.publicIp
		o = orm.NewOrm()
		_, err = o.Insert(&vps1)
		if err != nil {
			fmt.Println("vps插入失败", err)
			//rsp.Errno = utils.RECODE_DBERR
			//rsp.Errmsg = utils.RecodeText(rsp.Errno)
			return nil
		}
		o = orm.NewOrm()
		qs = o.QueryTable("vps")
		err = qs.Filter("usable_nodes__gt", 0).One(&vps)
		if err != nil {
			return err
		}
	}
	fmt.Printf("%d-%s-%s\n", vps.UsableNodes, order.Coin, vps.IpAddress)

	err = NewNode(vps.Id, order.Coin, vps.IpAddress, ssh_password)
	if err != nil {
		return err
	}

	node := models.Node{}
	node.Coin = order.Coin
	node.User = order.User
	node.Vps = &vps
	node.OrderNode = &order
	node.Port = 20000 - vps.UsableNodes*2
	o = orm.NewOrm()
	_, err = o.Insert(&node)
	if err != nil {
		fmt.Println("vps插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	fmt.Println("inser node finish")

	o = orm.NewOrm()
	vps.UsableNodes = vps.UsableNodes - 1
	o.Update(&vps)
	if err != nil {
		fmt.Println("vps update失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	fmt.Println("update vps finish")

	return nil
}

func TestProduct() error {
	product := models.Product{}
	product.Name = "masternode service"
	product.Period = "1 day"
	product.Amount = 1

	//插入数据到数据库中
	o := orm.NewOrm()
	_, err := o.Insert(&product)
	if err != nil {
		fmt.Println("Coin插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	product = models.Product{}
	product.Name = "masternode service"
	product.Period = "1 month"
	product.Amount = 25
	//插入数据到数据库中
	_, err = o.Insert(&product)
	if err != nil {
		fmt.Println("Coin插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	product = models.Product{}
	product.Name = "masternode service"
	product.Period = "6 month"
	product.Amount = 130
	//插入数据到数据库中
	_, err = o.Insert(&product)
	if err != nil {
		fmt.Println("Coin插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	product = models.Product{}
	product.Name = "masternode service"
	product.Period = "1 year"
	product.Amount = 208
	//插入数据到数据库中
	_, err = o.Insert(&product)
	if err != nil {
		fmt.Println("Coin插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	product = models.Product{}
	product.Name = "masternode service"
	product.Period = "3 year"
	product.Amount = 560
	//插入数据到数据库中
	_, err = o.Insert(&product)
	if err != nil {
		fmt.Println("Coin插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	return nil
}

func TestCoin() error {
	coin := models.Coin{}
	coin.Name = "vircle"
	coin.Status = "Enabled"
	coin.Path = ".vircle"
	coin.Conf = "vircle.conf"
	coin.Docker = "vpub/vircle:0.1"
	//插入数据到数据库中
	o := orm.NewOrm()
	_, err := o.Insert(&coin)
	if err != nil {
		fmt.Println("Coin插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	coin = models.Coin{}
	coin.Name = "dash"
	coin.Status = "Enabled"
	coin.Path = ".dashcore"
	coin.Conf = "dash.conf"
	coin.Docker = "vpub/dash:0.1"
	//插入数据到数据库中
	_, err = o.Insert(&coin)
	if err != nil {
		fmt.Println("Coin插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}

	coin = models.Coin{}
	coin.Name = "zcore"
	coin.Status = "Disabled"
	coin.Path = ".zcore"
	coin.Conf = "zcore.conf"
	coin.Docker = "vpub/zcore:0.1"
	//插入数据到数据库中
	_, err = o.Insert(&coin)
	if err != nil {
		fmt.Println("Coin插入失败", err)
		//rsp.Errno = utils.RECODE_DBERR
		//rsp.Errmsg = utils.RecodeText(rsp.Errno)
		return nil
	}
	return nil
}

func UploadFile(ipAddress string, username string, password string, localFile string, remoteDir string) error {
	log.Info(ipAddress, username, password, localFile, remoteDir)
	client, err := ssh.NewClient(ipAddress, "22", username, password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer client.Close()
	err = client.Upload(localFile, remoteDir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func WriteConf(confName string, port int) error {
	log.Info("conf start")
	cfg := ini.Empty()

	cfg.Section("").Key("rpcuser").SetValue(rpc_user)
	cfg.Section("").Key("rpcpassword").SetValue(rpc_password)
	cfg.Section("").Key("rpcallowip").SetValue("1.2.3.4/0.0.0.0")
	cfg.Section("").Key("rpcbind").SetValue("0.0.0.0")
	cfg.Section("").Key("rpcport").SetValue(strconv.Itoa(port))
	cfg.Section("").Key("port").SetValue(strconv.Itoa(port + 1))

	cfg.SaveTo(confName)
	log.Info("conf2 stop")

	return nil
}

func (e *Vps) processNewNode(userId string, orderId int64) error {
	log.Info("Processing new node start:", orderId)

	var order models.OrderNode
	o := orm.NewOrm()
	qs := o.QueryTable("order_node")
	err := qs.Filter("id", orderId).One(&order)
	if err != nil {
		e.pubErrMsg(userId, "newnode", utils.RECODE_NODATA, err.Error(), topic_newnode_fail)
		return err
	}

	var nvpsInfo VpsInfo
	var vpsInfo *VpsInfo
	var vps models.Vps
	nvps := models.Vps{}
	o = orm.NewOrm()
	qs = o.QueryTable("vps")
	err = qs.Filter("usable_nodes__gt", 0).One(&vps)
	retrys := 0
	if err != nil {
		fmt.Println("vps node 查询失败", err)
		for { //循环
			retrys++
			fmt.Println("vps retrys:", retrys)
			if retrys >= 20 {
				e.pubErrMsg(userId, "newnode", utils.TIMEOUT_VPS, err.Error(), topic_newnode_fail)
				return nil
			}
			vpsInfo, err = NewVps("", "", "", 1, nvpsInfo)
			if err != nil {
				time.Sleep(5 * time.Second)
				continue
			}
			fmt.Println("**********instanceId:", vpsInfo.volumeId)
			nvpsInfo.allocationId = vpsInfo.allocationId
			nvpsInfo.allocationState = vpsInfo.allocationState
			nvpsInfo.instanceId = vpsInfo.instanceId
			nvpsInfo.publicIp = vpsInfo.publicIp
			nvpsInfo.regionName = vpsInfo.regionName
			nvpsInfo.volumeId = vpsInfo.volumeId
			nvpsInfo.volumeState = vpsInfo.volumeState
			break
		}
		fmt.Println("success new vps:", vpsInfo.instanceId)

		log.Info(vpsInfo)
		nvps.AllocateId = vpsInfo.allocationId
		nvps.InstanceId = vpsInfo.instanceId
		nvps.VolumeId = vpsInfo.volumeId
		nvps.ProviderName = "amazon"
		nvps.Cores = 1
		nvps.Memory = 1
		nvps.KeyPairName = "vcl-keypair"
		nvps.MaxNodes = 5
		nvps.UsableNodes = 5
		nvps.SecurityGroupName = "vcl-mngroup"
		nvps.RegionName = vpsInfo.regionName
		nvps.IpAddress = vpsInfo.publicIp
		o = orm.NewOrm()
		_, err = o.Insert(&nvps)
		if err != nil {
			e.pubErrMsg(userId, "newnode", utils.RECODE_SERVERERR, err.Error(), topic_newnode_fail)
			return nil
		}
		o = orm.NewOrm()
		qs = o.QueryTable("vps")
		err = qs.Filter("usable_nodes__gt", 0).One(&vps)
		if err != nil {
			e.pubErrMsg(userId, "newnode", utils.RECODE_SERVERERR, err.Error(), topic_newnode_fail)
			return err
		}
	}
	fmt.Printf("%d-%s-%s\n", vps.UsableNodes, order.Coin, vps.IpAddress)

	retrys = 0
	for { //循环
		retrys++
		if retrys >= 5 {
			e.pubErrMsg(userId, "newnode", utils.TIMEOUT_VOLUME, "", topic_newnode_fail)
			return nil
		}
		err = NewNode(vps.Id, order.Coin, vps.IpAddress, ssh_password)
		fmt.Println("node retrys:", retrys)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}
	fmt.Println("success new node")

	rpcPort, _, err := getRpcPort(vps.Id)
	if err != nil {
		fmt.Println("get port 插入失败", err)
		e.pubErrMsg(userId, "newnode", utils.RECODE_INSERTERR, err.Error(), topic_newnode_fail)
		return nil
	}

	node := models.Node{}
	node.Coin = order.Coin
	node.User = order.User
	node.Vps = &vps
	node.OrderNode = &order
	node.Port = rpcPort
	o = orm.NewOrm()
	nodeId, err := o.Insert(&node)
	if err != nil {
		fmt.Println("node 插入失败", err)
		e.pubErrMsg(userId, "newnode", utils.RECODE_INSERTERR, err.Error(), topic_newnode_fail)
		return nil
	}
	fmt.Println("success insert node")

	o = orm.NewOrm()
	vps.UsableNodes = vps.UsableNodes - 1
	_, err = o.Update(&vps)
	if err != nil {
		fmt.Println("vps update失败", err)
		e.pubErrMsg(userId, "newnode", utils.RECODE_UPDATEERR, err.Error(), topic_newnode_fail)
		return nil
	}
	fmt.Println("update vps finish")

	o = orm.NewOrm()
	order.Status = models.ORDER_STATUS_COMPLETE
	_, err = o.Update(&order)
	if err != nil {
		fmt.Println("order update失败", err)
		e.pubErrMsg(userId, "newnode", utils.RECODE_UPDATEERR, err.Error(), topic_newnode_fail)
		return nil
	}
	fmt.Println("update order status finish")

	e.pubMsg(userId, topic_newnode_success, orderId)
	fmt.Println("success order:", nodeId)
	log.Info("Processing new node start:", orderId)
	return nil
}

func (e *Vps) processDelNode(userId string, nodeId int64) error {
	log.Info("Processing del node start:", nodeId)

	var node models.Node
	o := orm.NewOrm()
	qs := o.QueryTable("node")
	err := qs.Filter("id", nodeId).One(&node)
	if err != nil {
		e.pubErrMsg(userId, "delnode", utils.RECODE_NODATA, err.Error(), topic_newnode_fail)
		return err
	}

	var vps models.Vps
	o = orm.NewOrm()
	qs = o.QueryTable("vps")
	err = qs.Filter("id", node.Vps.Id).One(&vps)
	if err != nil {
		e.pubErrMsg(userId, "delnode", utils.RECODE_NODATA, err.Error(), topic_newnode_fail)
		return err
	}

	var order models.OrderNode
	o = orm.NewOrm()
	qs = o.QueryTable("order_node")
	err = qs.Filter("id", node.OrderNode.Id).One(&order)
	if err != nil {
		e.pubErrMsg(userId, "delnode", utils.RECODE_NODATA, err.Error(), topic_newnode_fail)
		return err
	}

	client := gossh.New(vps.IpAddress, "root")
	client.SetPassword(ssh_password)
	//client.DebugWriter = MakeLogger("DEBUG")
	client.InfoWriter = MakeLogger("INFO ")
	client.ErrorWriter = MakeLogger("ERROR")

	defer client.Close()

	port := node.Port
	log.Info("del volumn")
	volumeName := "mn" + strconv.Itoa(port)
	fmt.Println(volumeName)

	cmd := "docker stop  `docker ps -aq --filter name=" + volumeName + "`"
	sshCmd(client, cmd)
	cmd = "docker rm  `docker ps -aq --filter name=" + volumeName + "`"
	sshCmd(client, cmd)
	cmd = "docker volume rm " + volumeName
	sshCmd(client, cmd)

	o = orm.NewOrm()
	fmt.Println("del node:", node.Id, node.Vps.Id)
	_, err = o.Delete(&node)
	if err != nil {
		fmt.Println("node del", err)
		e.pubErrMsg(userId, "delnode", utils.RECODE_UPDATEERR, "", topic_delnode_fail)
		return nil
	}
	fmt.Println("update vps finish")

	o = orm.NewOrm()
	vps.UsableNodes = vps.UsableNodes + 1
	_, err = o.Update(&vps)
	if err != nil {
		fmt.Println("vps update失败", err)
		e.pubErrMsg(userId, "delnode", utils.RECODE_UPDATEERR, "", topic_delnode_fail)
		return nil
	}
	fmt.Println("update vps finish")

	o = orm.NewOrm()
	order.Status = models.ORDER_STATUS_CANCELED
	_, err = o.Update(&order)
	if err != nil {
		fmt.Println("order update失败", err)
		e.pubErrMsg(userId, "delnode", utils.RECODE_UPDATEERR, "", topic_delnode_fail)
		return nil
	}
	fmt.Println("update order status finish")

	isDel, err := e.delVps(userId, vps.Id)
	if err != nil {
		fmt.Println("del vps失败", err)
		e.pubErrMsg(userId, "delnode", utils.RECODE_UPDATEERR, "", topic_delnode_fail)
	}

	if isDel == true {
		o = orm.NewOrm()
		fmt.Println("del vps is:", vps.Id)
		_, err = o.Delete(&vps)
		if err != nil {
			fmt.Println("del vps record失败", err)
			e.pubErrMsg(userId, "delnode", utils.RECODE_UPDATEERR, "", topic_newnode_fail)
			return nil
		}
	}

	e.pubMsg(userId, topic_delnode_success, nodeId)
	fmt.Println("success del order:", nodeId)

	return nil
}

func (e *Vps) delVps(userId string, vpsId int) (bool, error) {
	fmt.Println("start del vps:", vpsId)
	var nodes []models.Node
	o := orm.NewOrm()
	qs := o.QueryTable("node")
	nums, err := qs.Filter("vps_id", vpsId).All(&nodes)
	if err != nil {
		e.pubErrMsg(userId, "delnode", utils.RECODE_NODATA, err.Error(), topic_newnode_fail)
		return false, err
	}
	fmt.Println("del nums:", nums)
	if nums == 0 {
		var vps models.Vps
		o = orm.NewOrm()
		qs = o.QueryTable("vps")
		err = qs.Filter("id", vpsId).One(&vps)
		if err != nil {
			e.pubErrMsg(userId, "delnode", utils.RECODE_NODATA, err.Error(), topic_newnode_fail)
			return false, err
		}

		c, err := uec2.NewEc2Client("", "test-account")
		if err != nil {
			fmt.Println("new ec2 error:", err)
			e.pubErrMsg(userId, "delnode", utils.RECODE_NODATA, err.Error(), topic_newnode_fail)
			return false, err
		}
		log.Info("client-vps")

		_, err = c.TerminateInstance(vps.InstanceId)
		if err != nil {
			e.pubErrMsg(userId, "delnode", utils.RECODE_NODATA, err.Error(), topic_newnode_fail)
			return false, err
		}

		retrys := 0
		for { //循环
			retrys++
			if retrys >= 30 {
				fmt.Println("get instance retrys:", retrys)
				e.pubErrMsg(userId, "delnode", utils.TIMEOUT_VPS, err.Error(), topic_newnode_fail)
				return false, err
			}
			result, err := c.GetDescribeInstance([]string{vps.InstanceId})
			fmt.Println("get instance retrys:", retrys)
			if err == nil {
				state := aws.StringValue(result.Reservations[0].Instances[0].State.Name)
				fmt.Println("vps state:", state)
				if state == "terminated" {
					break
				}
			}
			time.Sleep(time.Second)
		}

		log.Info("volumeId:", vps.VolumeId)
		_, err = c.DeleteVolumes(vps.VolumeId)
		if err != nil {
			e.pubErrMsg(userId, "delnode", utils.RECODE_NODATA, err.Error(), topic_newnode_fail)
			return false, err
		}

		_, err = c.ReleaseAddresss(vps.AllocateId)
		if err != nil {
			e.pubErrMsg(userId, "delnode", utils.RECODE_NODATA, err.Error(), topic_newnode_fail)
			return false, err
		}

		return true, nil
	}
	return false, nil
}

func getRpcPort(vpsId int) (int, int, error) {
	rpcport := port_from
	port := rpcport + 1

	fmt.Println("port start:", vpsId)
	var nodes []models.Node
	o := orm.NewOrm()
	qs := o.QueryTable("node")
	nums, err := qs.Filter("vps_id", vpsId).All(&nodes)
	if err != nil {
		return 0, 0, err
	}
	fmt.Println("nums:", nums)
	if nums == 0 {
		return rpcport, port, nil
	}

	for i := port_from; i < port_to; i = i + 2 {
		if portExist(i, &nodes) == false {
			return i, i + 1, nil
		}
	}
	return 0, 0, errors.New("portfull")
}

func portExist(port int, nodes *[]models.Node) bool {
	for _, node := range *nodes {
		if node.Port == port {
			return true
		}
	}
	return false
}

func sshCmd(client *gossh.Client, cmd string) (*gossh.Result, error) {
	retrys := 0
	for { //循环
		retrys++
		fmt.Println("ssh retrys:", retrys, cmd)
		if retrys >= 5 {
			return nil, errors.New("timeout")
		}
		result, err := client.Execute(cmd)
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(time.Second)
			continue
		}
		fmt.Println(result.Stdout())
		return result, nil
	}
}

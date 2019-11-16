module vircle.io/mnhost/api

go 1.13

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/micro/go-log v0.1.0 // indirect
	github.com/micro/go-plugins v1.5.1 // indirect
	github.com/uber/jaeger-client-go v2.20.1+incompatible // indirect
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	vircle.io/mnhost/common v0.0.0-00010101000000-000000000000 // indirect
	vircle.io/mnhost/config v0.0.0-00010101000000-000000000000 // indirect
	vircle.io/mnhost/interface/out/user v0.0.0-00010101000000-000000000000
)

replace github.com/gogo/protobuf v0.0.0-20190410021324-65acae22fc9 => github.com/gogo/protobuf v0.0.0-20190723190241-65acae22fc9d

replace vircle.io/mnhost/common => /root/mygo/src/vircle.io/mnhost/common

replace vircle.io/mnhost/go-plugins/broker/kafka => /root/mygo/src/vircle.io/mnhost/go-plugins/broker/kafka

replace vircle.io/mnhost/utils => /root/mygo/src/vircle.io/mnhost/utils

replace vircle.io/mnhost/model => /root/mygo/src/vircle.io/mnhost/model

replace vircle.io/mnhost/config => /root/mygo/src/vircle.io/mnhost/config

replace vircle.io/mnhost/interface/proto/vps => /root/mygo/src/vircle.io/mnhost/interface/proto/vps

replace vircle.io/mnhost/interface/out/nodemsg => /root/mygo/src/vircle.io/mnhost/interface/out/nodemsg

replace vircle.io/mnhost/interface/out/user => /root/mygo/src/vircle.io/mnhost/interface/out/user

replace vircle.io/mnhost/interface/out/vessel => /root/mygo/src/vircle.io/mnhost/interface/out/vessel

replace vircle.io/mnhost/interface/out/log => /root/mygo/src/vircle.io/mnhost/interface/out/log

replace vircle.io/mnhost/interface/out/vps => /root/mygo/src/vircle.io/mnhost/interface/out/vps

replace github.com/John-Tonny/micro/vps/amazon => /root/mygo/src/github.com/John-Tonny/micro/vps/amazon

replace github.com/golang/lint => golang.org/x/lint v0.0.0-20190930215403-16217165b5de

replace github.com/testcontainers/testcontainer-go => github.com/testcontainers/testcontainers-go v0.0.8
